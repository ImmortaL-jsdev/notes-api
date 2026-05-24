package worker

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func StartExportWorker(ctx context.Context, rdb *redis.Client) {
	log.Println("Export worker started")

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
		}
	}
}
