```markdown
# Notes API

REST API для управления заметками с аутентификацией, авторизацией и асинхронным экспортом. Построен на Go с использованием чистой архитектуры (handler → service → repository).

## Стек технологий

- **Go 1.25+**
- **PostgreSQL** — хранение заметок и пользователей
- **Redis** — очередь задач для асинхронного экспорта
- **JWT** — аутентификация и авторизация
- **Docker / docker-compose** — локальное окружение
- **Prometheus** — сбор метрик
- **GitHub Actions** — CI/CD (линтер, тесты, зеркалирование в GitLab)
- **Testcontainers** — интеграционные тесты с реальной БД

## Требования

- Go 1.25+
- Docker и docker-compose

## Быстрый старт

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/ImmortaL-jsdev/notes-api.git
   cd notes-api
   ```

2. Запустите PostgreSQL и Redis:
   ```bash
   docker compose up -d
   ```

3. Примените миграции (создание таблиц):
   ```bash
   docker compose exec -T postgres psql -U notes_user -d notes_db < migrations/000001_create_notes_table.up.sql
   docker compose exec -T postgres psql -U notes_user -d notes_db < migrations/000002_create_users.up.sql
   docker compose exec -T postgres psql -U notes_user -d notes_db < migrations/000003_add_user_id_to_notes.up.sql
   ```

4. Запустите сервер:
   ```bash
   go run cmd/server/main.go
   ```
   Сервер будет доступен на `http://localhost:8080`.

## Переменные окружения

| Переменная     | Описание                | Значение по умолчанию |
|----------------|-------------------------|------------------------|
| `DB_HOST`      | Хост PostgreSQL         | `localhost`            |
| `DB_PORT`      | Порт PostgreSQL         | `5432`                 |
| `DB_USER`      | Пользователь PostgreSQL | `notes_user`           |
| `DB_PASSWORD`  | Пароль PostgreSQL       | `notes_pass`           |
| `DB_NAME`      | Имя базы данных         | `notes_db`             |
| `JWT_SECRET`   | Секрет для JWT-токенов  | `supersecret`          |
| `REDIS_ADDR`   | Адрес Redis             | `localhost:6379`       |
| `PORT`         | Порт HTTP-сервера       | `8080`                 |

## API Endpoints

### Аутентификация (публичные)

- **Регистрация**  
  `POST /api/register`  
  Тело: `{"email":"user@example.com","password":"secret"}`

- **Вход**  
  `POST /api/login`  
  Тело: `{"email":"user@example.com","password":"secret"}`  
  Ответ: `{"access_token":"...", "refresh_token":"..."}`

### Заметки (защищённые, требуют заголовок `Authorization: Bearer <access_token>`)

- **Получить все заметки**  
  `GET /notes`

- **Создать заметку**  
  `POST /notes`  
  Тело: `{"title":"Заголовок","content":"Содержимое"}`

- **Получить заметку по ID**  
  `GET /notes/{id}`

- **Обновить заметку**  
  `PUT /notes/{id}`  
  Тело: `{"title":"Новый заголовок","content":"Новое содержимое"}`

- **Удалить заметку**  
  `DELETE /notes/{id}`

- **Массовое создание**  
  `POST /notes/bulk`  
  Тело: `[{"title":"Первая","content":"..."},{"title":"Вторая","content":"..."}]`

- **Долгая операция (демонстрация таймаутов)**  
  `GET /notes/process`

- **Экспорт заметок (асинхронный)**  
  `POST /notes/export`  
  Задача ставится в очередь Redis, результат сохраняется в папку `exports/`.

### Мониторинг

- **Метрики Prometheus**  
  `GET /metrics`

## Тестирование

### Юнит-тесты (сервис с моками)
```bash
make test-unit
```

### Интеграционные тесты (репозиторий с testcontainers)
```bash
make test-integration
```

### Все тесты
```bash
make test
```

### Линтер
```bash
make lint
```

## Структура проекта

```
notes-api/
├── cmd/server/            # точка входа
├── internal/
│   ├── auth/              # JWT-генерация и валидация
│   ├── errors/            # кастомные ошибки
│   ├── handlers/          # HTTP-обработчики
│   ├── middleware/        # JWT, метрики, логирование, восстановление
│   ├── models/            # структуры Note, User, TokenPair
│   ├── redis/             # клиент Redis
│   ├── repository/        # работа с PostgreSQL
│   ├── service/           # бизнес-логика
│   └── worker/            # фоновый воркер экспорта
├── migrations/            # SQL-миграции
├── .github/workflows/     # CI/CD
├── docker-compose.yml     # локальное окружение
├── Makefile
└── go.mod
```

## CI/CD

При каждом пуше в `main` запускаются:
- Линтер (`golangci-lint`)
- Юнит-тесты
- Интеграционные тесты (с testcontainers)

Код автоматически зеркалируется в GitLab.

## Лицензия

MIT
```
