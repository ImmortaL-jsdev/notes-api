package worker

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

			notes, err := store.GetAllForUser(ctx, userID)
			if err != nil {
				log.Printf("Failed to get notes for user %s: %v", userID, err)
				continue
			}

			fileName := filepath.Join("exports", fmt.Sprintf("%s_%d.txt", userID, time.Now().Unix()))

			file, err := os.Create(fileName)

			if err != nil {
				log.Printf("Failed to create file %s: %v", fileName, err)
				continue
			}

			for _, note := range notes {
				line := fmt.Sprintf("ID: %s\nTitle : %s\nContent: %s\nCreated: %s\n\n", note.ID, note.Title, note.Content, note.CreatedAt.Format(time.RFC3339))

				if _, err := file.WriteString(line); err != nil {
					log.Printf("Failed to write note %s: %v", note.ID, err)
				}
			}

			if err := file.Close(); err != nil {
				log.Printf("Failed to close file: %v", err)
			}
			log.Printf("Exported %d notes for user %s to %s", len(notes), userID, fileName)
		}
	}
}
