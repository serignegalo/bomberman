package game

import (
	"fmt"
	"math/rand"
	"strings"
)

type PLAYER struct {
	UserName     string     `json:"username"`
	Name         string     `json:"name"`
	LivesCount   int        `json:"livesCount"`
	Position     Position   `json:"position"`
	BeforeMove   Position   `json:"beforeMove"`
	Speed        int        `json:"speed"`
	Bomb         int        `json:"bomb"`
	Flame        int        `json:"flame"`
	Score        int        `json:"score"`
	Number       string     `json:"number"`
	BombPosition []Position `json:"bombPosition"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

const (
	PLAYERS_COUNT        = 4
	EMPTY_COUNT          = 30
	WIDTH                = 11
	HEIGHT               = 11
	LIVE_COUNT           = 3
	GAME_INIT            = "session_creation"
	GAME_RUNNING         = "playing_game"
	GAME_END             = "game_over"
	WAITING_PLAYERS      = "waiting_players"
	STARTING_GAME        = "starting_game"
	BLOCK_BRICK          = "b"
	BLOCK_WALL           = "#"
	BLOCK_FLAME_GIFT     = "gf"
	BLOCK_BOMB           = "bomb"
	BLOCK_BOMB_GIFT      = "gb"
	BLOCK_SPEED_GIFT     = "gs"
	BLOCK_EXPLOSION      = "explosion"
	BLOCK_HIDE_FLAME     = "F"
	BLOCK_HIDE_BOMB      = "B"
	BLOCK_HIDE_SPEED     = "S"
	EXPLOSION_HIDE_FLAME = "ef"
	EXPLOSION_HIDE_BOMB  = "eb"
	EXPLOSION_HIDE_SPEED = "es"
	BLOCK_PLAYER         = "player"
	BLOCK_EMPTY          = " "
	BLOCK_COIN           = "c"
	NBR_SPEED            = 2
	NBR_BOMB             = 2
	NBR_FLAME            = 2
)

func InitMap(width, height int) [][]string {
	grid := make([][]string, height)
	for i := range grid {
		grid[i] = make([]string, width)
	}
	return grid
}
func GenFixedWalls(grid [][]string) [][]string {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if (i+1)%2 == 1 && (j+1)%2 == 1 {
				grid[i][j] = BLOCK_WALL
			} else {
				grid[i][j] = BLOCK_EMPTY
			}
		}
	}
	return grid
}
func GenContour(grid [][]string) [][]string {
	for i := 0; i < len(grid[0]); i++ {
		grid[0][i] = BLOCK_WALL
		grid[len(grid)-1][i] = BLOCK_WALL
	}
	for i := 0; i < len(grid[0]); i++ {
		grid[i][0] = BLOCK_WALL
		grid[i][len(grid)-1] = BLOCK_WALL
	}
	return grid
}
func GetNumberWalls(grid [][]string) int {
	count := 0
	for _, row := range grid {
		for _, block := range row {
			if block == BLOCK_WALL {
				count++
			}
		}
	}
	return count
}
func MakeCoin(grid [][]string, tag string) [][]string {
	// Coin 1
	grid[1][1] = tag
	grid[1][2] = tag
	grid[1][3] = tag
	grid[2][1] = tag
	grid[3][1] = tag
	// Coin 2
	grid[1][WIDTH-2] = tag
	grid[1][WIDTH-3] = tag
	grid[1][WIDTH-4] = tag
	grid[2][WIDTH-2] = tag
	grid[3][WIDTH-2] = tag
	// Coin 3
	grid[HEIGHT-2][1] = tag
	grid[HEIGHT-2][2] = tag
	grid[HEIGHT-2][3] = tag
	grid[HEIGHT-3][1] = tag
	grid[HEIGHT-4][1] = tag
	// Coin 4
	grid[HEIGHT-2][WIDTH-2] = tag
	grid[HEIGHT-2][WIDTH-3] = tag
	grid[HEIGHT-2][WIDTH-4] = tag
	grid[HEIGHT-3][WIDTH-2] = tag
	grid[HEIGHT-4][WIDTH-2] = tag
	return grid
}
func GenBricks(grid [][]string) [][]string {
	nbrWalls := GetNumberWalls(grid)
	nbrBricks := (WIDTH * HEIGHT) - (nbrWalls + EMPTY_COUNT)
	var flatern []string
	for _, row := range grid {
		flatern = append(flatern, row...)
	}
	var emptyIndices []int
	for i, val := range flatern {
		if val == BLOCK_EMPTY {
			emptyIndices = append(emptyIndices, i)
		}
	}
	nbreBrick := RandomInt(nbrBricks-5, nbrBricks)
	compteur := 0
	for compteur < nbreBrick {
		randomIndice := RandomInt(0, len(emptyIndices)-1)
		flatern[emptyIndices[randomIndice]] = BLOCK_BRICK
		emptyIndices = removeElement(emptyIndices, randomIndice)
		compteur++
	}
	indice := 0
	for i, row := range grid {
		for j := range row {
			grid[i][j] = flatern[indice]
			indice++
		}
	}
	// Reconstitution du tableau en 2 dimension
	return grid
}
func GenPowers(grid [][]string) [][]string {
	var flatern []string
	for _, row := range grid {
		flatern = append(flatern, row...)
	}
	var brickIndices []int
	for i, val := range flatern {
		if val == BLOCK_BRICK {
			brickIndices = append(brickIndices, i)
		}
	}
	compteur := 0
	for compteur < NBR_BOMB {
		randomIndice := RandomInt(0, len(brickIndices)-1)
		flatern[brickIndices[randomIndice]] = BLOCK_HIDE_BOMB
		brickIndices = removeElement(brickIndices, randomIndice)
		compteur++
	}

	compteur = 0
	for compteur < NBR_FLAME {
		randomIndice := RandomInt(0, len(brickIndices)-1)
		flatern[brickIndices[randomIndice]] = BLOCK_HIDE_FLAME
		brickIndices = removeElement(brickIndices, randomIndice)
		compteur++
	}
	compteur = 0
	for compteur < NBR_SPEED {
		randomIndice := RandomInt(0, len(brickIndices)-1)
		flatern[brickIndices[randomIndice]] = BLOCK_HIDE_SPEED
		brickIndices = removeElement(brickIndices, randomIndice)
		compteur++
	}
	indice := 0
	for i, row := range grid {
		for j := range row {
			grid[i][j] = flatern[indice]
			indice++
		}
	}
	return grid
}
func PrintMap(grid [][]string) {
	for i := 0; i < len(grid); i++ {
		fmt.Println(strings.Join(grid[i], ""))
	}
}
func GenMap() [][]string {
	grid := InitMap(WIDTH, HEIGHT)
	grid = GenFixedWalls(grid)
	grid = GenContour(grid)
	grid = MakeCoin(grid, BLOCK_COIN)
	grid = GenBricks(grid)
	grid = GenPowers(grid)
	grid = MakeCoin(grid, BLOCK_EMPTY)
	return grid
}

// Fonction pour choisir un entier aléatoire dans la plage [min, max]
func RandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
func removeElement(arr []int, index int) []int {
	// Vérifiez si l'index est valide
	if index < 0 || index >= len(arr) {
		return arr
	}
	result := make([]int, 0, len(arr)-1)
	result = append(result, arr[:index]...)
	result = append(result, arr[index+1:]...)
	return result
}
func GeneratePosition(playersCount int, plateau [][]string) Position {
	var position = Position{X: 0, Y: 0}
	if playersCount == 1 {
		position.Y = 1
		position.X = 1
	} else if playersCount == 2 {
		position.X = WIDTH - 2
		position.Y = 1
	} else if playersCount == 3 {
		position.X = 1
		position.Y = HEIGHT - 2
	} else if playersCount == 4 {
		position.X = WIDTH - 2
		position.Y = HEIGHT - 2
	}
	return position
}

func ValidateUserName(userName string, players []PLAYER, state string) string {
	if userName == "" {
		return "name must not be empty"
	}
	if len(players) >= PLAYERS_COUNT {
		return "The number of player is full. Try in few minutes"
	}
	for _, player := range players {
		if player.UserName == userName {
			return "There are a player who have this name. Please choose unique username"
		}
	}
	if state == STARTING_GAME || state == GAME_RUNNING {
		return "There are a starting or running party. Please wait a few minute"
	}
	return ""
}

func CanMove(pos Position, grid [][]string, direction string) bool {
	if direction == "left" {
		pos.X--
		if pos.X < 0 {
			return false
		}
	} else if direction == "right" {
		pos.X++
		if pos.X > WIDTH-1 {
			return false
		}
	} else if direction == "up" {
		pos.Y--
		if pos.Y < 0 {
			return false
		}
	} else if direction == "down" {
		pos.Y++
		if pos.Y > HEIGHT-1 {
			return false
		}
	}
	if grid[pos.Y][pos.X] == BLOCK_WALL ||
		grid[pos.Y][pos.X] == BLOCK_BRICK ||
		grid[pos.Y][pos.X] == BLOCK_BOMB ||
		grid[pos.Y][pos.X] == BLOCK_EXPLOSION ||
		grid[pos.Y][pos.X] == BLOCK_HIDE_BOMB ||
		grid[pos.Y][pos.X] == BLOCK_HIDE_FLAME ||
		grid[pos.Y][pos.X] == BLOCK_HIDE_SPEED ||
		grid[pos.Y][pos.X] == "player1" ||
		grid[pos.Y][pos.X] == "player2" ||
		grid[pos.Y][pos.X] == "player3" ||
		grid[pos.Y][pos.X] == "player4" {
		return false
	}
	return true
}

// pour déplacer un joueur
func MovePlayer(mapGame [][]string, player PLAYER, sense string) ([][]string, PLAYER) {
	if CanMove(player.Position, mapGame, sense) {
		mapGame = CleanPlayer(player.Position, mapGame)
		player.Position = UpdatePosition(player.Position, sense)
		player = GetGift(mapGame, player)
		mapGame = Place(player.Position, mapGame, player.Number)
	}
	return mapGame, player
}
func PlaceFlame(grid [][]string, pos Position) [][]string {

	if pos.Y < 0 || pos.Y >= HEIGHT {
		return grid
	} else if pos.X < 0 || pos.X >= WIDTH {
		return grid
	} else if grid[pos.Y][pos.X] == BLOCK_WALL {
		return grid
	} else if grid[pos.Y][pos.X] == BLOCK_HIDE_BOMB ||
		grid[pos.Y][pos.X] == BLOCK_BOMB_GIFT {
		grid[pos.Y][pos.X] = EXPLOSION_HIDE_BOMB
		return grid
	} else if grid[pos.Y][pos.X] == BLOCK_HIDE_SPEED ||
		grid[pos.Y][pos.X] == BLOCK_SPEED_GIFT {
		grid[pos.Y][pos.X] = EXPLOSION_HIDE_SPEED
		return grid
	} else if grid[pos.Y][pos.X] == BLOCK_HIDE_FLAME ||
		grid[pos.Y][pos.X] == BLOCK_FLAME_GIFT {
		grid[pos.Y][pos.X] = EXPLOSION_HIDE_FLAME
		return grid
	} else {
		grid[pos.Y][pos.X] = BLOCK_EXPLOSION
		return grid
	}
}
func GetGift(grid [][]string, player PLAYER) PLAYER {
	if grid[player.Position.Y][player.Position.X] == BLOCK_BOMB_GIFT {
		player.Bomb++
	} else if grid[player.Position.Y][player.Position.X] == BLOCK_FLAME_GIFT {
		player.Flame++
	} else if grid[player.Position.Y][player.Position.X] == BLOCK_SPEED_GIFT {
		if player.Speed == 50 {
			player.Speed = 20
		} else {
			player.Speed = 50
		}
	}
	return player
}
func Place(pos Position, grid [][]string, block string) [][]string {
	grid[pos.Y][pos.X] = block
	return grid
}
func CleanPlayer(pos Position, grid [][]string) [][]string {
	if grid[pos.Y][pos.X] == "player1" ||
		grid[pos.Y][pos.X] == "player2" ||
		grid[pos.Y][pos.X] == "player3" ||
		grid[pos.Y][pos.X] == "player4" {
		grid[pos.Y][pos.X] = BLOCK_EMPTY
	}
	return grid
}
func GetBlock(pos Position, grid [][]string) string {
	return grid[pos.Y][pos.X]
}

func GetPlayer(players []PLAYER, username string) (PLAYER, int, error) {
	for i, player := range players {
		if player.UserName == username {
			return player, i, nil
		}
	}
	return PLAYER{}, 0, fmt.Errorf("player does not exist")
}

func UpdatePosition(pos Position, direction string) Position {
	//fmt.Println("avant--------updated pos:", pos)
	switch {
	case direction == "left":
		pos.X--
	case direction == "right":
		pos.X++
	case direction == "up":
		pos.Y--
	case direction == "down":
		pos.Y++
	}
	//fmt.Println("apres------updated pos:", pos)
	return pos
}
func PlaceAllFlame(mapGame [][]string, player PLAYER) [][]string {
	for _, bombPosition := range player.BombPosition {

		//to place bomb on position player
		mapGame = PlaceFlame(mapGame, bombPosition)
		//to place bomb on left
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X - i,
				Y: bombPosition.Y,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = PlaceFlame(mapGame, pos)
		}

		//to place bomb on right
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X + i,
				Y: bombPosition.Y,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = PlaceFlame(mapGame, pos)
		}
		//to place bomb on down
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X,
				Y: bombPosition.Y - i,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = PlaceFlame(mapGame, pos)
		}
		//to place bomb on top
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X,
				Y: bombPosition.Y + i,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = PlaceFlame(mapGame, pos)
		}
	}
	return mapGame
}
func CleanFlame(grid [][]string, pos Position) [][]string {

	if pos.Y < 0 || pos.Y >= HEIGHT {
		return grid
	} else if pos.X < 0 || pos.X >= WIDTH {
		return grid
	} else if grid[pos.Y][pos.X] == BLOCK_WALL {
		return grid
	} else if grid[pos.Y][pos.X] == EXPLOSION_HIDE_BOMB {
		grid[pos.Y][pos.X] = BLOCK_BOMB_GIFT
		return grid
	} else if grid[pos.Y][pos.X] == EXPLOSION_HIDE_SPEED {
		grid[pos.Y][pos.X] = BLOCK_SPEED_GIFT
		return grid
	} else if grid[pos.Y][pos.X] == EXPLOSION_HIDE_FLAME {
		grid[pos.Y][pos.X] = BLOCK_FLAME_GIFT
		return grid
	} else {
		grid[pos.Y][pos.X] = BLOCK_EMPTY
		return grid
	}
}
func LivesCount(mapGame [][]string, posFlame Position, players []PLAYER) []PLAYER {
	for i, player := range players {
		if player.LivesCount != 0 {
			mapGame = Place(player.Position, mapGame, player.Number)
		}else{
			mapGame = Place(player.Position, mapGame, BLOCK_EMPTY)
		}
		if posFlame.X == player.Position.X && posFlame.Y == player.Position.Y {
			if player.LivesCount > 0 {
				player.LivesCount--
			}else{
				mapGame = Place(player.Position, mapGame, BLOCK_EMPTY)
			}
			players[i] = player
		}
	}
	return players
}
func CleanAllFlame(mapGame [][]string, player PLAYER, players []PLAYER) ([][]string, []PLAYER) {
	for _, bombPosition := range player.BombPosition {

		//to place bomb on position player
		mapGame = CleanFlame(mapGame, bombPosition)
		players = LivesCount(mapGame, bombPosition, players)
		//to place bomb on left
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X - i,
				Y: bombPosition.Y,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = CleanFlame(mapGame, pos)
			players = LivesCount(mapGame, pos, players)
		}

		//to place bomb on right
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X + i,
				Y: bombPosition.Y,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = CleanFlame(mapGame, pos)
			players = LivesCount(mapGame, pos, players)
		}
		//to place bomb on down
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X,
				Y: bombPosition.Y - i,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = CleanFlame(mapGame, pos)
			players = LivesCount(mapGame, pos, players)
		}
		for i := 1; i <= player.Flame; i++ {
			pos := Position{
				X: bombPosition.X,
				Y: bombPosition.Y + i,
			}
			if mapGame[pos.Y][pos.X] == BLOCK_WALL {
				break
			}
			mapGame = CleanFlame(mapGame, pos)
			players = LivesCount(mapGame, pos, players)
		}
	}
	return mapGame, players
}

func PoseAllBomb(mapGame [][]string, player PLAYER) PLAYER {
	player.BombPosition = append(player.BombPosition, player.Position)
	posedBomb := 1
	indice := 1
	for posedBomb < player.Bomb {
		pos := Position{
			X: player.Position.X + indice,
			Y: player.Position.Y,
		}
		if pos.Y < 0 || pos.Y >= HEIGHT {
			break
		} else if pos.X < 0 || pos.X >= WIDTH {
			break
		} else if mapGame[pos.Y][pos.X] != BLOCK_EMPTY {
			break
		}
		player.BombPosition = append(player.BombPosition, pos)
		posedBomb++
		indice++
	}
	indice = 1
	for posedBomb < player.Bomb {
		pos := Position{
			X: player.Position.X,
			Y: player.Position.Y + indice,
		}
		if pos.Y < 0 || pos.Y >= HEIGHT {
			break
		} else if pos.X < 0 || pos.X >= WIDTH {
			break
		} else if mapGame[pos.Y][pos.X] != BLOCK_EMPTY {
			break
		}
		player.BombPosition = append(player.BombPosition, pos)
		posedBomb++
		indice++
	}
	indice = 1
	for posedBomb < player.Bomb {
		pos := Position{
			X: player.Position.X,
			Y: player.Position.Y - indice,
		}
		if pos.Y < 0 || pos.Y >= HEIGHT {
			break
		} else if pos.X < 0 || pos.X >= WIDTH {
			break
		} else if mapGame[pos.Y][pos.X] != BLOCK_EMPTY {
			break
		}
		player.BombPosition = append(player.BombPosition, pos)
		posedBomb++
		indice++
	}
	indice = 1
	for posedBomb < player.Bomb {
		pos := Position{
			X: player.Position.X - indice,
			Y: player.Position.Y,
		}
		if pos.Y < 0 || pos.Y >= HEIGHT {
			break
		} else if pos.X < 0 || pos.X >= WIDTH {
			break
		} else if mapGame[pos.Y][pos.X] != BLOCK_EMPTY {
			break
		}
		player.BombPosition = append(player.BombPosition, pos)
		posedBomb++
		indice++
	}
	return player
}
