package main

import (
	"log"
	"net/http"
	"todo-api/databases"
	"todo-api/storage"
)

func main() {
	dbConfig := databases.Config{
		Host:     "localhost",
		Port:     5434,
		User:     "postgres",
		Password: "example",
		DBName:   "todo-db",
		SSlMode:  "disable", // Use "require" in production
	}
	migrationsPath := "./migrations"

	// Connect to database
	db, err := databases.NewDBConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	if err := databases.RunMigrations(db, migrationsPath); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

	store := storage.NewMemoryStore()

	server := NewTodoServer(store)

	log.Print("Server is running on port 5000")
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}

}
