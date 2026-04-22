package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/store"
	"github.com/gorilla/mux"
)

func TestGetAllNotes(t *testing.T) {
	store := store.NewMemoryStore()

	store.Create(models.Note{Title: "test", Content: "content"})

	handler := NewNoteHandler(store)

	req := httptest.NewRequest("GET", "/notes", nil)

	rr := httptest.NewRecorder()

	handler.GetAll(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}

	var notes []models.Note

	err := json.Unmarshal(rr.Body.Bytes(), &notes)

	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(notes))
	}

	if notes[0].Title != "test" {
		t.Errorf("Wrong title")
	}
}

func TestCreateNotes(t *testing.T) {
	store := store.NewMemoryStore()

	handler := NewNoteHandler(store)

	newNote := models.Note{Title: "Buy milk", Content: "2 liters"}

	jsonData, err := json.Marshal(newNote)

	if err != nil {
		t.Fatalf("marshal error : %v", err)
	}

	req := httptest.NewRequest("POST", "/notes", bytes.NewReader(jsonData))

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("status = %d, want 201", rr.Code)
	}

	var createdNote models.Note

	err = json.Unmarshal(rr.Body.Bytes(), &createdNote)

	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if createdNote.ID == "" {
		t.Errorf("ID is empty")
	}

	if createdNote.CreatedAt.IsZero() {
		t.Errorf("createdAt is zero")
	}

	if createdNote.Title != newNote.Title {
		t.Errorf("title = %s, want %s", createdNote.Title, newNote.Title)
	}

	all := store.GetAll()
	if len(all) != 1 {
		t.Errorf("expected 1 note, got %d", len(all))
	}
}

func TestGetNotes(t *testing.T) {
	store := store.NewMemoryStore()

	existingNote, err := store.Create(models.Note{Title: "test", Content: "content"})

	if err != nil {
		t.Fatalf("failed to create store")
	}

	handler := NewNoteHandler(store)

	url := "/notes/" + existingNote.ID

	req := httptest.NewRequest("GET", url, nil)

	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/notes/{id}", handler.GetByID).Methods("GET")
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}

	var getNotes models.Note

	err = json.Unmarshal(rr.Body.Bytes(), &getNotes)

	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if getNotes.ID != existingNote.ID {
		t.Errorf("id = %s, want %s", getNotes.ID, existingNote.ID)
	}

	if getNotes.Content != existingNote.Content {
		t.Errorf("content = %s, want %s", getNotes.Content, existingNote.Content)
	}

	if getNotes.Title != existingNote.Title {
		t.Errorf("title = %s, want %s", getNotes.Title, existingNote.Title)
	}

	if !getNotes.CreatedAt.Equal(existingNote.CreatedAt) {
		t.Errorf("createdAt = %v, want %v", getNotes.CreatedAt, existingNote.CreatedAt)
	}

}

func TestGetNoteNotFound(t *testing.T) {
	store := store.NewMemoryStore()
	handler := NewNoteHandler(store)

	req := httptest.NewRequest("GET", "/notes/123", nil)
	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()
	handler.GetByID(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func TestUpdateNote(t *testing.T) {
	store := store.NewMemoryStore()

	originalNote := models.Note{Title: "test", Content: "content"}

	created, err := store.Create(originalNote)

	if err != nil {
		t.Fatalf("failed to create note : %v", err)
	}

	handler := NewNoteHandler(store)

	updatedNote := models.Note{Title: "new title", Content: "new content"}

	jsonData, err := json.Marshal(updatedNote)

	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/notes/{id}", handler.Update).Methods("PUT")

	req := httptest.NewRequest("PUT", "/notes/"+created.ID, bytes.NewReader(jsonData))

	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rr.Code)
	}

	var got models.Note

	err = json.Unmarshal(rr.Body.Bytes(), &got)

	if err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}

	if got.ID != created.ID {
		t.Errorf("id = %s, want %s", got.ID, created.ID)
	}

	if got.Content != updatedNote.Content {
		t.Errorf("content = %s, want %s", got.Content, updatedNote.Content)
	}

	if got.Title != updatedNote.Title {
		t.Errorf("title = %s, want %s", got.Title, updatedNote.Title)
	}

	if !got.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("createdAt = %v, want %v", got.CreatedAt, created.CreatedAt)
	}

}

func TestUpdateNoteNotFound(t *testing.T) {
	store := store.NewMemoryStore()
	handler := NewNoteHandler(store)
	r := mux.NewRouter()
	r.HandleFunc("/notes/{id}", handler.Update).Methods("PUT")

	updatedNote := models.Note{Title: "new title", Content: "new content"}
	jsonData, _ := json.Marshal(updatedNote)

	req := httptest.NewRequest("PUT", "/notes/123", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}

func TestDeleteNote(t *testing.T) {
	store := store.NewMemoryStore()

	originalNote := models.Note{Title: "test", Content: "content"}

	created, err := store.Create(originalNote)

	if err != nil {
		t.Fatalf("failed to create note : %v", err)
	}

	handler := NewNoteHandler(store)

	r := mux.NewRouter()

	r.HandleFunc("/notes/{id}", handler.Delete).Methods("DELETE")

	req := httptest.NewRequest("DELETE", "/notes/"+created.ID, nil)

	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", rr.Code)
	}

	_, ok := store.GetByID(created.ID)

	if ok {
		t.Errorf("note still exists after delete:")
	}

}

func TestDeleteNoteNotFound(t *testing.T) {

	store := store.NewMemoryStore()
	handler := NewNoteHandler(store)

	req := httptest.NewRequest("DELETE", "/notes/123", nil)
	req.Header.Set("X-API-Key", "secret123")

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	r.HandleFunc("/notes/{id}", handler.Delete).Methods("DELETE")

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rr.Code)
	}
}
