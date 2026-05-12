package repository_test

import (
	"context"
	"testing"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresStore_Create_Integration(t *testing.T) {

	ctx := context.Background()
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
	)
	if err != nil {
		t.Fatalf("не удалось запустить контейнер: %v", err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("не удалось остановить контейнер: %v", err)
		}
	}()

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Fatalf("не удалось создать временный пул: %v", err)
	}
	defer pool.Close()

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("не удалось создать таблицу notes: %v", err)
	}
	pool.Close()

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("не удалось создать store: %v", err)
	}
	defer store.Close()

	note := models.Note{Title: "Integration Test", Content: "Hello from container"}
	created, err := store.Create(ctx, note)
	if err != nil {
		t.Fatalf("ошибка при создании заметки: %v", err)
	}

	if created.ID == "" {
		t.Error("ID не должен быть пустым")
	}
	if created.Title != note.Title {
		t.Errorf("ожидался Title %q, получен %q", note.Title, created.Title)
	}
	if created.Content != note.Content {
		t.Errorf("ожидался Content %q, получен %q", note.Content, created.Content)
	}
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt не должен быть нулевым")
	}
}

func TestPostgresStore_GetAll_Integration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"))

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("не удалось остановить контейнер: %v", err)
		}
	}()

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("не удалось создать таблицу notes: %v", err)
	}
	pool.Close()

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("не удалось создать store: %v", err)
	}
	defer store.Close()

	/////////////////////////////////////////////////////////

	notesToCreate := []models.Note{{Title: "Первая", Content: "Раз"}, {Title: "Вторая", Content: "Два"}}

	var createdNotes []models.Note

	for _, n := range notesToCreate {
		created, err := store.Create(ctx, n)
		if err != nil {
			t.Fatal(err)
		}
		createdNotes = append(createdNotes, created)
	}

	getAll, err := store.GetAll(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if len(getAll) != len(createdNotes) {
		t.Errorf("ожидалось %d заметок, получено %d", len(createdNotes), len(getAll))
	}

	createdMap := make(map[string]bool)
	for _, n := range createdNotes {
		createdMap[n.ID] = true
	}
	for _, n := range getAll {
		if !createdMap[n.ID] {
			t.Errorf("найден лишний ID: %s", n.ID)
		}
	}
}

func TestPostgresStore_GetByID_Integration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"))

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("не удалось остановить контейнер: %v", err)
		}
	}()

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("не удалось создать таблицу notes: %v", err)
	}
	pool.Close()

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("не удалось создать store: %v", err)
	}
	defer store.Close()

	originalNote := models.Note{Title: "GetByID тест", Content: "Содержимое"}

	created, err := store.Create(ctx, originalNote)

	if err != nil {
		t.Fatal(err)
	}

	fetched, err := store.GetByID(ctx, created.ID)

	if err != nil {
		t.Fatal(err)
	}

	if fetched.ID != created.ID {
		t.Errorf("ID не совпадает: %s vs %s", fetched.ID, created.ID)
	}
	if fetched.Title != originalNote.Title {
		t.Errorf("Title не совпадает: %s vs %s", fetched.Title, originalNote.Title)
	}
	if fetched.Content != originalNote.Content {
		t.Errorf("Content не совпадает: %s vs %s", fetched.Content, originalNote.Content)
	}
	if !fetched.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("CreatedAt не совпадает: %v vs %v", fetched.CreatedAt, created.CreatedAt)
	}

	_, err = store.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего ID, но её нет")
	}

}

func TestPostgresStore_Update_Integration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"))

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("не удалось остановить контейнер: %v", err)
		}
	}()

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("не удалось создать таблицу notes: %v", err)
	}
	pool.Close()

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("не удалось создать store: %v", err)
	}
	defer store.Close()

	originalNote := models.Note{Title: "Старый", Content: "Старое"}

	created, err := store.Create(ctx, originalNote)
	if err != nil {
		t.Fatal(err)
	}

	updatedData := models.Note{Title: "Новый", Content: "Новое"}
	updated, err := store.Update(ctx, created.ID, updatedData)
	if err != nil {
		t.Fatal(err)
	}

	if updated.ID != created.ID {
		t.Errorf("ID изменился: %s vs %s", updated.ID, created.ID)
	}
	if updated.Title != updatedData.Title {
		t.Errorf("Title не обновился: %s vs %s", updated.Title, updatedData.Title)
	}
	if updated.Content != updatedData.Content {
		t.Errorf("Content не обновился: %s vs %s", updated.Content, updatedData.Content)
	}
	if !updated.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("CreatedAt изменился: %v vs %v", updated.CreatedAt, created.CreatedAt)
	}

	_, err = store.Update(ctx, "00000000-0000-0000-0000-000000000000", updatedData)
	if err == nil {
		t.Error("ожидалась ошибка при обновлении несуществующей заметки")

	}
}

func TestPostgresStore_Delete_Integration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"))

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("не удалось остановить контейнер: %v", err)
		}
	}()

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("не удалось создать таблицу notes: %v", err)
	}
	pool.Close()

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("не удалось создать store: %v", err)
	}
	defer store.Close()

	note := models.Note{Title: "Удалить", Content: "Содержимое"}
	created, err := store.Create(ctx, note)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Delete(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.GetByID(ctx, created.ID)
	if err == nil {
		t.Error("заметка не удалена, GetByID вернул результат")
	}

	err = store.Delete(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("ожидалась ошибка при удалении несуществующей заметки")
	}
}

func TestPostgresStore_CreateMany_Integration(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx, testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"))

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("не удалось остановить контейнер: %v", err)
		}
	}()

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

	if err != nil {
		t.Fatalf("не удалось получить строку подключения: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)

	createTableSQL := `CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);`
	_, err = pool.Exec(ctx, createTableSQL)
	if err != nil {
		t.Fatalf("не удалось создать таблицу notes: %v", err)
	}
	pool.Close()

	store, err := repository.NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("не удалось создать store: %v", err)
	}
	defer store.Close()

	validNotes := []models.Note{
		{Title: "Много 1", Content: "Содержимое 1"},
		{Title: "Много 2", Content: "Содержимое 2"},
	}
	created, err := store.CreateMany(ctx, validNotes)
	if err != nil {
		t.Fatalf("ошибка при массовом создании: %v", err)
	}
	if len(created) != 2 {
		t.Errorf("длина результата %d, ожидалось 2", len(created))
	}

	all, err := store.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 2 {
		t.Errorf("в таблице %d записей, ожидалось 2", len(all))
	}

	invalidNotes := []models.Note{
		{Title: "Корректная", Content: "good"},
		{Title: "", Content: "пустой заголовок"},
	}
	_, err = store.CreateMany(ctx, invalidNotes)
	if err == nil {
		t.Error("ожидалась ошибка при создании с пустым заголовком")
	}

	allAfter, err := store.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(allAfter) != 2 {
		t.Errorf("после ошибочной транзакции в таблице %d записей, ожидалось 2", len(allAfter))
	}
}
