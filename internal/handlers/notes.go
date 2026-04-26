package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"
	"github.com/gorilla/mux"
)

type NoteHandler struct {
	store *repository.PostgresStore
}

func NewNoteHandler(store *repository.PostgresStore) *NoteHandler {
	return &NoteHandler{store: store}
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.store.GetAll(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}
	created, err := h.store.Create(r.Context(), note)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	note, err := h.store.GetByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedNote models.Note
	err := json.NewDecoder(r.Body).Decode(&updatedNote)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}
	note, err := h.store.Update(r.Context(), id, updatedNote)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.store.Delete(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "note not found"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	err := json.NewDecoder(r.Body).Decode(&notes)
	if err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if len(notes) == 0 {
		http.Error(w, `{"error":"empty list"}`, http.StatusBadRequest)
		return
	}
	created, err := h.store.CreateMany(r.Context(), notes)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
