package service

import (
	"context"
	"fmt"

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

func (s *NoteService) GetAll(ctx context.Context) ([]models.Note, error) {
	return s.repo.GetAll(ctx)
}

func (s *NoteService) GetByID(ctx context.Context, id string) (models.Note, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NoteService) Update(ctx context.Context, id string, note models.Note) (models.Note, error) {
	return s.repo.Update(ctx, id, note)
}

func (s *NoteService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
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
