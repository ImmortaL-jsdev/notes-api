# Этап 1: Сборка (Builder)
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей и качаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь остальной код проекта
COPY . .

# Собираем бинарник (CGO_ENABLED=0 нужен для работы на Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/notes-api ./cmd/server/main.go

# Этап 2: Финальный легковесный образ
FROM alpine:latest

WORKDIR /app

# Копируем готовый бинарник из первого этапа
COPY --from=builder /app/notes-api .

# Открываем порт
EXPOSE 8080

# Команда для запуска
CMD ["/app/notes-api"]