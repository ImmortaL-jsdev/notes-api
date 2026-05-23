package service

import (
	"context"
	"fmt"
	"time"

	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
)

type NoteRepository interface {
	Create(ctx context.Context, note models.Note) (models.Note, error)
	GetAll(ctx context.Context) ([]models.Note, error)
	GetByID(ctx context.Context, id string) (models.Note, error)
	Update(ctx context.Context, id string, note models.Note) (models.Note, error)
	Delete(ctx context.Context, id string) error
	CreateMany(ctx context.Context, notes []models.Note) ([]models.Note, error)

	GetAllForUser(ctx context.Context, userID string) ([]models.Note, error)
	CreateForUser(ctx context.Context, userID string, note models.Note) (models.Note, error)
	GetByIDForUser(ctx context.Context, userID, noteID string) (models.Note, error)
	UpdateForUser(ctx context.Context, userID, noteID string, note models.Note) (models.Note, error)
	DeleteForUser(ctx context.Context, userID, noteID string) error
	CreateManyForUser(ctx context.Context, userID string, notes []models.Note) ([]models.Note, error)
}

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) Create(ctx context.Context, note models.Note) (models.Note, error) {
	if note.Title == "" {
		return models.Note{}, &myerrors.ValidationError{Message: "title cannot be empty"}
	}

	created, err := s.repo.Create(ctx, note)

	if err != nil {
		return models.Note{}, fmt.Errorf("failed to create note: %w", err)
	}

	return created, nil
}

func (s *NoteService) CreateForUser(ctx context.Context, userID string, note models.Note) (models.Note, error) {
	if note.Title == "" {
		return models.Note{}, &myerrors.ValidationError{Message: "title cannot be empty"}
	}
	return s.repo.CreateForUser(ctx, userID, note)
}

func (s *NoteService) GetAll(ctx context.Context) ([]models.Note, error) {
	return s.repo.GetAll(ctx)
}

func (s *NoteService) GetAllForUser(ctx context.Context, userID string) ([]models.Note, error) {
	return s.repo.GetAllForUser(ctx, userID)
}

func (s *NoteService) GetByID(ctx context.Context, id string) (models.Note, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NoteService) GetByIDForUser(ctx context.Context, userID, noteID string) (models.Note, error) {
	return s.repo.GetByIDForUser(ctx, userID, noteID)
}

func (s *NoteService) Update(ctx context.Context, id string, note models.Note) (models.Note, error) {
	return s.repo.Update(ctx, id, note)
}

func (s *NoteService) UpdateForUser(ctx context.Context, userID string, id string, note models.Note) (models.Note, error) {
	return s.repo.UpdateForUser(ctx, userID, id, note)
}

func (s *NoteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *NoteService) DeleteForUser(ctx context.Context, userID, noteID string) error {
	return s.repo.DeleteForUser(ctx, userID, noteID)
}

func (s *NoteService) CreateMany(ctx context.Context, notes []models.Note) ([]models.Note, error) {
	for _, note := range notes {
		if note.Title == "" {

			return nil, &myerrors.ValidationError{Message: "title cannot be empty"}
		}
	}
	created, err := s.repo.CreateMany(ctx, notes)
	if err != nil {
		return nil, fmt.Errorf("failed to create many notes: %w", err)
	}
	return created, nil
}

func (s *NoteService) CreateManyForUser(ctx context.Context, userID string, notes []models.Note) ([]models.Note, error) {
	for _, note := range notes {
		if note.Title == "" {
			return nil, &myerrors.ValidationError{Message: "title cannot be empty"}
		}
	}
	return s.repo.CreateManyForUser(ctx, userID, notes)
}

func (s *NoteService) Process(ctx context.Context) error {
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
