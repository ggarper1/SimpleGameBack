package main

import (
	"log"
	"net/http"

	"ggarper1/SimpleGameBack/src/routes"
)

func main() {
	manager := routes.NewMatchCreator()

	http.HandleFunc("/ws", manager.RecieveConnection)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
