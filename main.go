package main

import (
	"fmt"
	"main/handlers"
	"net/http"
)

const PORT = ":4044"

func main() {
	static := http.FileServer(http.Dir("./public/"))
	http.HandleFunc("/", handlers.AppHandler)
	http.HandleFunc("/validate", handlers.ValidateUserNameHandler)
	http.HandleFunc("/ws", handlers.ConnectionHandler)
	http.Handle("/public/", http.StripPrefix("/public/", static))
	fmt.Println("Server running on http://localhost" + PORT)
	//go handlers.HandleMessages()
	http.ListenAndServe(PORT, nil)
}
