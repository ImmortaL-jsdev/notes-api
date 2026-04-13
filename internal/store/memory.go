package store

import (
	"errors"
	"sync"
	"time"

	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/google/uuid"
)

type MemoryStore struct {
	notes map[string]models.Note
	mu    sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		notes: make(map[string]models.Note),
	}
}

func (s *MemoryStore) Create(note models.Note) (models.Note, error) {

	if note.Title == "" {
		return models.Note{}, errors.New("title cannot be empty")
	}

	id := uuid.New().String()

	note.ID = id
	note.CreatedAt = time.Now().UTC()

	s.mu.Lock()

	defer s.mu.Unlock()
	s.notes[id] = note

	return note, nil
}

func (s *MemoryStore) GetAll() []models.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var notesSlice []models.Note
	for _, note := range s.notes {
		notesSlice = append(notesSlice, note)
	}

	return notesSlice
}

func (s *MemoryStore) GetByID(id string) (models.Note, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	note, ok := s.notes[id]
	return note, ok
}
func (s *MemoryStore) Update(id string, updated models.Note) (models.Note, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	note, ok := s.notes[id]
	var updatedNote models.Note

	if ok {
		updatedNote = models.Note{
			ID:        note.ID,
			Title:     updated.Title,
			Content:   updated.Content,
			CreatedAt: note.CreatedAt,
		}
		s.notes[id] = updatedNote
		return updatedNote, true
	}
	return models.Note{}, false
}

func (s *MemoryStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.notes[id]
	if ok {
		delete(s.notes, id)
		return true
	}
	return false
}
