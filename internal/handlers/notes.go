package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/service"
	"github.com/gorilla/mux"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

type NoteHandler struct {
	service *service.NoteService
}

func NewNoteHandler(service *service.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	notes, err := h.service.GetAll(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	respondWithJSON(w, http.StatusOK, notes)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	var note models.Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	created, err := h.service.Create(r.Context(), note)
	if err != nil {
		var valErr *myerrors.ValidationError
		if errors.As(err, &valErr) {
			respondWithError(w, http.StatusBadRequest, valErr.Message)
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, created)
}

func (h *NoteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	note, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		var notFound *myerrors.NotFoundError
		if errors.As(err, &notFound) {
			respondWithError(w, http.StatusNotFound, notFound.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedNote models.Note
	if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	note, err := h.service.Update(r.Context(), id, updatedNote)
	if err != nil {
		var notFound *myerrors.NotFoundError
		if errors.As(err, &notFound) {
			respondWithError(w, http.StatusNotFound, notFound.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.service.Delete(r.Context(), id)
	if err != nil {
		var notFound *myerrors.NotFoundError
		if errors.As(err, &notFound) {
			respondWithError(w, http.StatusNotFound, notFound.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *NoteHandler) CreateBulk(w http.ResponseWriter, r *http.Request) {
	var notes []models.Note
	if err := json.NewDecoder(r.Body).Decode(&notes); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	if len(notes) == 0 {
		respondWithError(w, http.StatusBadRequest, "empty list")
		return
	}
	created, err := h.service.CreateMany(r.Context(), notes)
	if err != nil {
		var valErr *myerrors.ValidationError
		if errors.As(err, &valErr) {
			respondWithError(w, http.StatusBadRequest, valErr.Message)
		} else {
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, created)
}
