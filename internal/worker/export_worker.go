package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"
	"github.com/redis/go-redis/v9"
)

func StartExportWorker(ctx context.Context, rdb *redis.Client, store *repository.PostgresStore) {
	log.Println("Export worker started")

	if err := os.MkdirAll("exports", 0755); err != nil {
		log.Printf("Failed to create exports directory: %v", err)
		return

	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Export worker shutting down...")
			return
		default:
			result, err := rdb.BRPop(ctx, 0, "export-queue").Result()
			if err != nil {
				continue
			}
			userID := result[1]
			log.Printf("Processing export for user: %s", userID)

			const maxRetries = 3
			const retryDelay = 2 * time.Second

			var notes []models.Note
			var lastErr error

			for attempt := 1; attempt <= maxRetries; attempt++ {
				n, err := store.GetAllForUser(ctx, userID)
				if err == nil {
					notes = n
					lastErr = nil
					break
				}

				lastErr = err
				log.Printf("Attempt %d failed to user %s: %v", attempt, userID, err)

				if attempt < maxRetries {
					time.Sleep(retryDelay)
				}
			}

			if lastErr != nil {
				if err := rdb.LPush(ctx, "export-dlq", userID).Err(); err != nil {
					log.Printf("Failed to push to DLQ: %v", err)
				}
				log.Printf("Moved user %s to DLQ after %d attempts", userID, maxRetries)
				continue
			}

			fileName := filepath.Join("exports", fmt.Sprintf("%s_%d.txt", userID, time.Now().Unix()))

			file, err := os.Create(fileName)

			if err != nil {
				lastErr = fmt.Errorf("create file : %w", err)
				log.Printf("Attempt failed to create file : %v", lastErr)
				time.Sleep(retryDelay)
				continue
			}

			for _, note := range notes {
				line := fmt.Sprintf("ID: %s\nTitle : %s\nContent: %s\nCreated: %s\n\n", note.ID, note.Title, note.Content, note.CreatedAt.Format(time.RFC3339))

				if _, err := file.WriteString(line); err != nil {
					lastErr = fmt.Errorf("write note: %w", err)
					log.Printf("Failed to write note %s: %v", note.ID, err)
					break
				}
			}

			if err := file.Close(); err != nil {
				log.Printf("Failed to close file: %v", err)
			}
			log.Printf("Exported %d notes for user %s to %s", len(notes), userID, fileName)

		}
	}
}
