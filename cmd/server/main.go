package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ImmortaL-jsdev/notes-api/internal/handlers"
	"github.com/ImmortaL-jsdev/notes-api/internal/middleware"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"
	"github.com/ImmortaL-jsdev/notes-api/internal/service"
	"github.com/gorilla/mux"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "notes_user"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "notes_pass"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "notes_db"
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer store.Close()

	svc := service.NewNoteService(store)

	handler := handlers.NewNoteHandler(svc)

	r := mux.NewRouter()

	r.HandleFunc("/notes", handler.GetAll).Methods("GET")
	r.HandleFunc("/notes", handler.Create).Methods("POST")
	r.HandleFunc("/notes/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/notes/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/notes/{id}", handler.Delete).Methods("DELETE")
	r.HandleFunc("/notes/bulk", handler.CreateBulk).Methods("POST")

	r.Use(middleware.RecoveryMiddleware, middleware.LoggingMiddleware, middleware.AuthMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
