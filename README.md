```markdown
# Notes API

REST API для управления заметками на Go. Поддерживает создание, чтение, обновление и удаление заметок. Данные хранятся в памяти (in-memory store). Реализована аутентификация через API-ключ, логирование запросов и восстановление после паник.

## Требования

- Go 1.24+

## Установка и запуск

```bash
git clone https://github.com/ImmortaL-jsdev/notes-api.git
cd notes-api
go mod tidy
go run cmd/server/main.go
```

Сервер запустится на `http://localhost:8080`.

## API эндпоинты

Все запросы требуют заголовок:
```
X-API-Key: secret123
```

### Получить все заметки
```bash
curl -H "X-API-Key: secret123" http://localhost:8080/notes
```

### Создать заметку
```bash
curl -X POST -H "Content-Type: application/json" -H "X-API-Key: secret123" -d '{"title":"Buy milk","content":"2 liters"}' http://localhost:8080/notes
```

### Получить заметку по ID
```bash
curl -H "X-API-Key: secret123" http://localhost:8080/notes/{id}
```

### Обновить заметку
```bash
curl -X PUT -H "Content-Type: application/json" -H "X-API-Key: secret123" -d '{"title":"New title","content":"New content"}' http://localhost:8080/notes/{id}
```

### Удалить заметку
```bash
curl -X DELETE -H "X-API-Key: secret123" http://localhost:8080/notes/{id}
```

## Тестирование

Запуск всех тестов:
```bash
go test -v ./...
```

## Структура проекта

```
notes-api/
├── cmd
├── internal
│ ├── handlers
│ ├── middleware
│ ├── models
│ └── store
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Использование Makefile (опционально)

```bash
make run     # запустить сервер
make test    # запустить тесты
make lint    # запустить линтер
make build   # скомпилировать бинарник
make clean   # удалить бинарник
```

## Лицензия

MIT
```