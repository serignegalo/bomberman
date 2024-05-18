package handlers

import (
	"html/template"
	"main/game"
	"net/http"
)
var players []game.PLAYER
var state = game.GAME_INIT
// Pour gerer la page app
func AppHandler(w http.ResponseWriter, r *http.Request) {

	temp, err := template.ParseFiles("templates/app.html")
	if err != nil {
		http.Error(w, "Internal server", http.StatusInternalServerError)
		return
	}
	err = temp.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal server", http.StatusInternalServerError)
	}
}
func ValidateUserNameHandler(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("username")
	name := r.URL.Query().Get("name")
	if error := game.ValidateUserName(userName, players, state); error != "" {
		http.Error(w, error, http.StatusBadRequest)
		//fmt.Println("error validation user: ", error)
		return
	}
	player := game.PLAYER{
					UserName: userName,
					Name:       name,
					LivesCount: game.LIVE_COUNT,
					Speed:      400,
					Flame:      1,
					Bomb:       1,
					Score:      0,
				}
	players = append(players, player)
	http.Error(w, "ok", http.StatusOK)
}
