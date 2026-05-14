package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	go func() {
		log.Printf("Server is starting on:", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exited gracefully")

}
