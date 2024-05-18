package handlers

import (
	"fmt"
	"log"
	"main/game"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var mutex sync.Mutex
var clients = make(map[*websocket.Conn]string)
var broadcast = make(chan DataAction)
var waiting_players = 10
var waiting_start = 2
var mapGame [][]string

type ChatMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type MessageResponse struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
}
type DataAction struct {
	Type     string `json:"type"`
	Sense    string `json:"sense"`
	Username string `json:"username"`
	Message string      `json:"message"`

}
type DataResponse struct {
	Map        [][]string    `json:"map"`
	Player     game.PLAYER   `json:"player"`
	ListPlayer []game.PLAYER `json:"listPlayer"`
}

func ConnectionHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username")
	// creating connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//http.Error(w, "Internal server", http.StatusInternalServerError)
		fmt.Println("error creation connection:", err)
		return
	}
	mutex.Lock()
	clients[conn] = userName
	mutex.Unlock()
	if len(players)== 1{
		waitingPlayersNotif()
		if waiting_players == 0 && len(players) <= 1 {
			NotCompleteNotif()
			return
		}else{
			startingGameNotif()
			gameStartedNotif()
		}
	}
	// Boucle pour écouter les messages des clients
	go func() {
		for {
			var action DataAction
			// Lecture du message JSON
			err := conn.ReadJSON(&action)
			if err != nil {
				return
			}
			if action.Type == "move" {
				player, indice, err := game.GetPlayer(players, action.Username)
				if err == nil && player.LivesCount > 0 {
					beforeMove := player.Position
					mapGame, player = game.MovePlayer(mapGame, player, action.Sense)
					player.BeforeMove =  beforeMove
					players[indice] = player
					//fmt.Println("pos:", player.Position)
					dataResponse := DataResponse{
						Map:    mapGame,
						Player: player,
					}
					message := MessageResponse{
						Type: "move",
						Data: dataResponse,
					}
					for client := range clients {
						err = client.WriteJSON(message)
						if err != nil {
							return
						}
					}
				}
			} else if action.Type == "pose-bomb" {
				player, indice, err := game.GetPlayer(players, action.Username)
				if err == nil && player.LivesCount > 0 {
					player = game.PoseAllBomb(mapGame, player)
					players[indice] = player
					for _, bombPosition := range player.BombPosition {
						mapGame = game.Place(bombPosition, mapGame, game.BLOCK_BOMB)
					}
					dataResponse := DataResponse{
						Map:    mapGame,
						Player: player,
					}
					message := MessageResponse{
						Type: "pose-bomb",
						Data: dataResponse,
					}
					for client := range clients {
						err = client.WriteJSON(message)
						if err != nil {
							return
						}
					}
				}
			} else if action.Type == "explose-bomb" {
				player, indice, err := game.GetPlayer(players, action.Username)
				if err == nil && player.LivesCount > 0 {
					mapGame = game.PlaceAllFlame(mapGame, player)
					dataResponse := DataResponse{
						Map:        mapGame,
						Player:     player,
						ListPlayer: players,
					}
					message := MessageResponse{
						Type: "explose-bomb",
						Data: dataResponse,
					}
					for client := range clients {
						err = client.WriteJSON(message)
						if err != nil {
							return
						}
					}
					time.Sleep(1 * time.Second)
					mapGame, players = game.CleanAllFlame(mapGame, player, players)
					if players[indice].UserName == player.UserName {
						player.LivesCount = players[indice].LivesCount
					}
					player.BombPosition = []game.Position{}
					players[indice] = player
					message.Type = "after-explosion"
					for client := range clients {
						err = client.WriteJSON(message)
						if err != nil {
							return
						}
						//fmt.Println("client username:", username)
					}
				}
			} else if action.Type == "message" {
				chatMsg := ChatMessage{Username: action.Username, Message: action.Message}

				//logic traitement message de discussion
				for client := range clients {

					msg := MessageResponse{Type: "message", Data: chatMsg}

					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("Error sending chat message: %v", err)
						client.Close()
						delete(clients, client)
						return
					}
				}
			}
		}
	}()
}
func HandleMessages() {
	// Boucle pour envoyer des messages

	for {
		action := <-broadcast
		fmt.Println(action)
		if action.Type == "move" {
			player, indice, err := game.GetPlayer(players, action.Username)
			if err == nil && player.LivesCount > 0 {
				mapGame, player = game.MovePlayer(mapGame, player, action.Sense)
				players[indice] = player
				//fmt.Println("pos:", player.Position)
				dataResponse := DataResponse{
					Map:    mapGame,
					Player: player,
				}
				message := MessageResponse{
					Type: "move",
					Data: dataResponse,
				}
				for client := range clients {
					err = client.WriteJSON(message)
					if err != nil {
						return
					}
				}
			}
		} else if action.Type == "pose-bomb" {
			player, indice, err := game.GetPlayer(players, action.Username)
			if err == nil && player.LivesCount > 0 {
				player = game.PoseAllBomb(mapGame, player)
				players[indice] = player
				for _, bombPosition := range player.BombPosition {
					mapGame = game.Place(bombPosition, mapGame, game.BLOCK_BOMB)
				}
				dataResponse := DataResponse{
					Map:    mapGame,
					Player: player,
				}
				message := MessageResponse{
					Type: "pose-bomb",
					Data: dataResponse,
				}
				for client := range clients {
					err = client.WriteJSON(message)
					if err != nil {
						return
					}
				}
			}
		} else if action.Type == "explose-bomb" {
			player, indice, err := game.GetPlayer(players, action.Username)
			if err == nil && player.LivesCount > 0 {
				mapGame = game.PlaceAllFlame(mapGame, player)
				dataResponse := DataResponse{
					Map:        mapGame,
					Player:     player,
					ListPlayer: players,
				}
				message := MessageResponse{
					Type: "explose-bomb",
					Data: dataResponse,
				}
				for client := range clients {
					err = client.WriteJSON(message)
					if err != nil {
						return
					}
				}
				time.Sleep(1 * time.Second)
				mapGame, players = game.CleanAllFlame(mapGame, player, players)
				if players[indice].UserName == player.UserName {
					player.LivesCount = players[indice].LivesCount
				}
				player.BombPosition = []game.Position{}
				players[indice] = player

				message.Type = "after-explosion"
				for client := range clients {
					err = client.WriteJSON(message)
					if err != nil {
						return
					}
					//fmt.Println("client username:", username)
				}
			}
		} else if action.Type == "message" {
			//logic traitement message de discussion
		}
	}
}
func waitingPlayersNotif() {
	for waiting_players > 0 && len(clients) < game.PLAYERS_COUNT {
		state = game.WAITING_PLAYERS
		msg := MessageResponse{Type: state, Data: waiting_players}
		// Envoi du message à tous les clients connectés
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				return
			}
		}
		time.Sleep(1 * time.Second)
		if len(clients) == 4 {
			return
		}
		waiting_players--
	}
}
func startingGameNotif() {
	for waiting_start > 0 {
		state = game.STARTING_GAME
		msg := MessageResponse{Type: state, Data: waiting_start}
		// Envoi du message à tous les clients connectés
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				return
			}
		}
		// Attente de 1 seconde avant de décrémenter
		time.Sleep(1 * time.Second)
		waiting_start--
	}
}
func gameStartedNotif() {
	if state == game.GAME_RUNNING {
		return
	}
	mapGame = game.GenMap()
	for i := range players {
		pos := game.GeneratePosition(i+1, mapGame)
		players[i].Position = pos
		players[i].BeforeMove = pos
		players[i].Number = "player" + strconv.Itoa(i+1)
		mapGame[pos.Y][pos.X] = "player" + strconv.Itoa(i+1)
	}
	i := 0
	// Envoi du message à tous les clients connectés
	for client := range clients {
		data := DataResponse{
			Map:    mapGame,
			Player: players[i],
		}
		msg := MessageResponse{Type: "running_game", Data: data}
		err := client.WriteJSON(msg)
		if err != nil {
			return
		}
		i++
	}
	state = game.GAME_RUNNING
}
func NotCompleteNotif() {
	msg := MessageResponse{Type: "note_complete_game", Data: "You must have at least 2 players"}
	// Envoi du message à tous les clients connectés
	for client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			return
		}
	}
	initGame()
}
func initGame() {
	clients = make(map[*websocket.Conn]string)
	broadcast = make(chan DataAction)
	waiting_players = 20
	waiting_start = 10
	players = []game.PLAYER{}
	state = game.GAME_INIT
	mapGame = [][]string{}
	//fmt.Println(players)
}
