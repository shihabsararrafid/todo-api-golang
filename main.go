package main

import (
	"log"
	"net/http"
	"todo-api/storage"
)

func main() {

	store := storage.NewMemoryStore()

	server := NewTodoServer(store)

	log.Print("Server is running on port 5000")
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}

}
