package main

import (
	"log"
	"net/http"

	"github.com/ImmortaL-jsdev/notes-api/internal/handlers"
	"github.com/ImmortaL-jsdev/notes-api/internal/middleware"
	"github.com/ImmortaL-jsdev/notes-api/internal/store"
	"github.com/gorilla/mux"
)

func main() {
	s := store.NewMemoryStore()

	h := handlers.NewNoteHandler(s)

	r := mux.NewRouter()

	r.HandleFunc("/notes", h.GetAll).Methods("GET")
	r.HandleFunc("/notes", h.Create).Methods("POST")
	r.HandleFunc("/notes/{id}", h.GetByID).Methods("GET")
	r.HandleFunc("/notes/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/notes/{id}", h.Delete).Methods("DELETE")

	r.Use(middleware.RecoveryMiddleware, middleware.LoggingMiddleware, middleware.AuthMiddleware)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
