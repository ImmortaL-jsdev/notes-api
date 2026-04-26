package main

import (
	"log"
	"net/http"

	"github.com/ImmortaL-jsdev/notes-api/internal/handlers"
	"github.com/ImmortaL-jsdev/notes-api/internal/middleware"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"
	"github.com/gorilla/mux"
)

func main() {

	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer store.Close()

	h := handlers.NewNoteHandler(store)

	r := mux.NewRouter()

	r.HandleFunc("/notes", h.GetAll).Methods("GET")
	r.HandleFunc("/notes", h.Create).Methods("POST")
	r.HandleFunc("/notes/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/notes/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/notes/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/notes/bulk", h.CreateBulk).Methods("POST")

	r.Use(middleware.RecoveryMiddleware, middleware.LoggingMiddleware, middleware.AuthMiddleware)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
