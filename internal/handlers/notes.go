package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/store"
	"github.com/gorilla/mux"
)

type NoteHandler struct {
	store *store.MemoryStore
}

func NewNoteHandler(store *store.MemoryStore) *NoteHandler {
	return &NoteHandler{
		store: store,
	}
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes) //nolint:errcheck
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"}) //nolint:errcheck
		return
	}

	createdNote, err := h.store.Create(note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}) //nolint:errcheck
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdNote) //nolint:errcheck
	}
}
func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	note, ok := h.store.GetByID(id)
	w.Header().Set("Content-Type", "application/json")
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note) //nolint:errcheck
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note not found"}) //nolint:errcheck
	}
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	var updatedNote models.Note
	vars := mux.Vars(r)
	id := vars["id"]

	err := json.NewDecoder(r.Body).Decode(&updatedNote)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"}) //nolint:errcheck
		return
	}

	note, ok := h.store.Update(id, updatedNote)

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note not found"}) //nolint:errcheck
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note) //nolint:errcheck
	}
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	ok := h.store.Delete(id)

	if !ok {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note not found"}) //nolint:errcheck
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
