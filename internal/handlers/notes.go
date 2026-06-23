package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/middleware"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/service"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
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
	service     *service.NoteService
	redisClient *redis.Client
}

func NewNoteHandler(service *service.NoteService, rdb *redis.Client) *NoteHandler {
	return &NoteHandler{
		service:     service,
		redisClient: rdb,
	}
}

func (h *NoteHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	notes, err := h.service.GetAllForUser(ctx, userID)

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

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	created, err := h.service.CreateForUser(ctx, userID, note)

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

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	note, err := h.service.GetByIDForUser(ctx, userID, id)

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

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	note, err := h.service.UpdateForUser(ctx, userID, id, updatedNote)
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

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.service.DeleteForUser(ctx, userID, id)

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

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	created, err := h.service.CreateManyForUser(ctx, userID, notes)

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

func (h *NoteHandler) Process(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err := h.service.Process(ctx)

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		respondWithError(w, http.StatusRequestTimeout, "request timeout")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "process completed!"})

}

func (h *NoteHandler) Export(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)

	if !ok || userID == "" {
		respondWithError(w, http.StatusUnauthorized, "missing user id")
		return
	}

	if err := h.redisClient.LPush(r.Context(), "export-queue", userID).Err(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to enqueue export")
		return
	}

	respondWithJSON(w, http.StatusAccepted, map[string]string{"status": "export queued"})
}
