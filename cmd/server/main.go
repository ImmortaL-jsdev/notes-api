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
	redis_client "github.com/ImmortaL-jsdev/notes-api/internal/redis"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"
	"github.com/ImmortaL-jsdev/notes-api/internal/service"
	"github.com/ImmortaL-jsdev/notes-api/internal/worker"
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

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		jwtSecret = "supersecret"
	}

	redisAddr := os.Getenv("REDIS_ADDR")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb, err := redis_client.NewClient(context.Background(), redisAddr)

	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("redis close error: %v", err)
		}
	}()

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	defer store.Close()

	ctxWorker, cancelWorker := context.WithCancel(context.Background())

	defer cancelWorker()

	go worker.StartExportWorker(ctxWorker, rdb, store)

	userStore, err := repository.NewUserStore(connString)

	if err != nil {
		log.Fatal("Failed to connect to user store", err)
	}
	defer userStore.Close()

	svc := service.NewNoteService(store)
	handler := handlers.NewNoteHandler(svc, rdb)

	authService := service.NewAuthService(userStore, []byte(jwtSecret))
	authHandler := handlers.NewAuthHandler(authService)

	r := mux.NewRouter()

	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	api := r.PathPrefix("/notes").Subrouter()
	api.Use(middleware.AuthMiddleware)

	api.HandleFunc("", handler.GetAll).Methods("GET")
	api.HandleFunc("", handler.Create).Methods("POST")
	api.HandleFunc("/{id}", handler.GetByID).Methods("GET")
	api.HandleFunc("/{id}", handler.Update).Methods("PUT")
	api.HandleFunc("/{id}", handler.Delete).Methods("DELETE")
	api.HandleFunc("/bulk", handler.CreateBulk).Methods("POST")
	api.HandleFunc("/process", handler.Process).Methods("GET")
	api.HandleFunc("/export", handler.Export).Methods("POST")

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
		log.Printf("Server is starting on: %v", port)
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
