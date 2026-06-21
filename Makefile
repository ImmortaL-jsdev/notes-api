.PHONY: run test test-unit test-integration lint build clean docker-up docker-down

# Запустить сервер
run:
	go run cmd/server/main.go

# Запустить все тесты (юнит + интеграционные)
test: test-unit test-integration

# Только юнит-тесты (быстро, без Docker)
test-unit:
	go test -v ./internal/service/...

# Только интеграционные тесты (нужен Docker и флаг для Ryuk)
test-integration:
	TESTCONTAINERS_RYUK_DISABLED=true go test -v ./internal/repository/...

# Проверка стиля и типичных ошибок
lint:
	golangci-lint run ./...

# Сборка бинарника
build:
	go build -o bin/notes-api cmd/server/main.go

# Удалить собранный бинарник
clean:
	rm -rf bin/

# Запустить PostgreSQL и Redis
docker-up:
	docker compose up -d

# Остановить PostgreSQL и Redis
docker-down:
	docker compose down