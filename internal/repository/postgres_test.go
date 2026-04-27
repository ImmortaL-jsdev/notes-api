package repository

import (
	"context"
	"testing"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
)

func TestPostgresStore_Create(t *testing.T) {

	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"

	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("Не удалось подключиться: %v", err)
	}
	defer store.Close()

	_, err = store.pool.Exec(context.Background(), "TRUNCATE notes")
	if err != nil {
		t.Fatalf("Не удалось очистить таблицу: %v", err)
	}

	original := models.Note{Title: "Тест", Content: "Содержание"}
	created, err := store.Create(context.Background(), original)
	if err != nil {
		t.Fatalf("Ошибка при создании: %v", err)
	}

	if created.ID == "" {
		t.Error("ID не должен быть пустым")
	}
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt не должен быть нулевым")
	}
	if created.Title != original.Title || created.Content != original.Content {
		t.Error("Заголовок или содержимое не совпадают")
	}

	fetched, err := store.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Не удалось прочитать запись: %v", err)
	}
	if fetched.Title != original.Title {
		t.Error("Заголовок прочитанной записи не совпадает")
	}
}

func TestPostgresStore_GetAll(t *testing.T) {
	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"

	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatalf("Не удалось подключиться: %v", err)
	}
	defer store.Close()

	_, err = store.pool.Exec(context.Background(), "TRUNCATE notes")
	if err != nil {
		t.Fatalf("Не удалось очистить таблицу: %v", err)
	}

	notesToCreate := []models.Note{{Title: "Первая", Content: "Раз"}, {Title: "Вторая", Content: "Два"}}

	var createdNote []models.Note

	ctx := context.Background()

	for _, note := range notesToCreate {
		created, err := store.Create(ctx, note)

		if err != nil {
			t.Fatal(err)
		}

		createdNote = append(createdNote, created)
	}

	all, err := store.GetAll(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(all) != len(createdNote) {
		t.Errorf("ожидалось %d заметок, получено %d", len(createdNote), len(all))
	}

	createdMap := make(map[string]bool)

	for _, note := range createdNote {
		createdMap[note.ID] = true
	}

	for _, note := range all {
		if !createdMap[note.ID] {
			t.Errorf("найден лишний ID: %s", note.ID)
		}
	}
}

func TestPostgresStore_GetByID(t *testing.T) {
	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"
	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	_, err = store.pool.Exec(context.Background(), "TRUNCATE notes")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	original := models.Note{Title: "GetByID тест", Content: "Содержимое"}
	created, err := store.Create(ctx, original)
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
	if fetched.Title != original.Title {
		t.Errorf("Title не совпадает: %s vs %s", fetched.Title, original.Title)
	}
	if fetched.Content != original.Content {
		t.Errorf("Content не совпадает: %s vs %s", fetched.Content, original.Content)
	}
	if !fetched.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("CreatedAt не совпадает: %v vs %v", fetched.CreatedAt, created.CreatedAt)
	}

	_, err = store.GetByID(ctx, "00000000-0000-0000-0000-000000000000")
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующего ID, но её нет")
	}
}
func TestPostgresStore_Update(t *testing.T) {
	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"
	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	_, err = store.pool.Exec(context.Background(), "TRUNCATE notes")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	original := models.Note{Title: "Старый", Content: "Старое"}
	created, err := store.Create(ctx, original)
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

func TestPostgresStore_Delete(t *testing.T) {
	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"
	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	_, err = store.pool.Exec(context.Background(), "TRUNCATE notes")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

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

func TestPostgresStore_CreateMany(t *testing.T) {
	connString := "postgres://notes_user:notes_pass@localhost:5432/notes_db?sslmode=disable"
	store, err := NewPostgresStore(connString)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	_, err = store.pool.Exec(context.Background(), "TRUNCATE notes")
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

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
