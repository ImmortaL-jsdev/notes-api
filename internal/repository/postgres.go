package repository

import (
	"context"

	"fmt"

	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(connString string) (*PostgresStore, error) {
	pool, err := pgxpool.New(context.Background(), connString)

	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func (s *PostgresStore) Create(ctx context.Context, note models.Note) (models.Note, error) {
	var created models.Note

	query := `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id, created_at`

	err := s.pool.QueryRow(ctx, query, note.Title, note.Content).Scan(&created.ID, &created.CreatedAt)

	if err != nil {
		return models.Note{}, fmt.Errorf("failed to insert note: %w", err)
	}
	created.Title = note.Title
	created.Content = note.Content
	return created, nil
}

func (s *PostgresStore) CreateForUser(ctx context.Context, userID string, note models.Note) (models.Note, error) {
	var created models.Note
	query := `INSERT INTO notes (title, content, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, user_id`
	err := s.pool.QueryRow(ctx, query, note.Title, note.Content, userID).Scan(&created.ID, &created.CreatedAt, &created.UserID)
	if err != nil {
		return models.Note{}, fmt.Errorf("failed to insert note: %w", err)
	}
	created.Title = note.Title
	created.Content = note.Content
	return created, nil
}

func (s *PostgresStore) GetAll(ctx context.Context) ([]models.Note, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, title, content, created_at FROM notes ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notes, nil
}
func (s *PostgresStore) GetAllForUser(ctx context.Context, userID string) ([]models.Note, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, title, content, created_at, user_id FROM notes WHERE user_id = $1 ORDER BY created_at DESC",
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UserID); err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return notes, nil
}

func (s *PostgresStore) GetByID(ctx context.Context, id string) (models.Note, error) {
	var note models.Note
	query := `SELECT id, title, content, created_at FROM notes WHERE id = $1`

	err := s.pool.QueryRow(ctx, query, id).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt)

	if err != nil {
		return models.Note{}, &myerrors.NotFoundError{Entity: "note", ID: id}
	}
	return note, nil
}

func (s *PostgresStore) GetByIDForUser(ctx context.Context, userID, noteID string) (models.Note, error) {
	var note models.Note
	query := `SELECT id, title, content, created_at, user_id FROM notes WHERE id = $1 AND user_id = $2`

	err := s.pool.QueryRow(ctx, query, noteID, userID).Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UserID)

	if err != nil {
		return models.Note{}, &myerrors.NotFoundError{Entity: "note", ID: noteID}
	}

	return note, nil
}

func (s *PostgresStore) Update(ctx context.Context, id string, note models.Note) (models.Note, error) {
	if _, err := s.GetByID(ctx, id); err != nil {
		return models.Note{}, err
	}

	query := `UPDATE notes SET title = $1, content = $2 WHERE id = $3 RETURNING id, title, content, created_at`
	var updatedNote models.Note
	err := s.pool.QueryRow(ctx, query, note.Title, note.Content, id).Scan(&updatedNote.ID, &updatedNote.Title, &updatedNote.Content, &updatedNote.CreatedAt)

	if err != nil {
		return models.Note{}, fmt.Errorf("failed to update note: %w", err)
	}
	return updatedNote, nil
}

func (s *PostgresStore) UpdateForUser(ctx context.Context, userID, noteID string, note models.Note) (models.Note, error) {
	if _, err := s.GetByIDForUser(ctx, userID, noteID); err != nil {
		return models.Note{}, err
	}

	query := `UPDATE notes SET title = $1, content = $2 WHERE id = $3 AND user_id = $4 RETURNING id, title, content, created_at, user_id`
	var updatedNote models.Note
	err := s.pool.QueryRow(ctx, query, note.Title, note.Content, noteID, userID).Scan(&updatedNote.ID, &updatedNote.Title, &updatedNote.Content, &updatedNote.CreatedAt, &updatedNote.UserID)

	if err != nil {
		return models.Note{}, fmt.Errorf("failed to update note: %w", err)
	}
	return updatedNote, nil
}

func (s *PostgresStore) Delete(ctx context.Context, id string) error {
	cmdTag, err := s.pool.Exec(ctx, "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return &myerrors.NotFoundError{Entity: "note", ID: id}
	}
	return nil
}

func (s *PostgresStore) DeleteForUser(ctx context.Context, userID, noteID string) error {
	cmdTag, err := s.pool.Exec(ctx, "DELETE FROM notes WHERE id = $1 AND user_id = $2", noteID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return &myerrors.NotFoundError{Entity: "note", ID: noteID}
	}
	return nil
}

func (s *PostgresStore) CreateMany(ctx context.Context, notes []models.Note) ([]models.Note, error) {
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	created := make([]models.Note, 0, len(notes))

	for _, note := range notes {

		if note.Title == "" {
			txErr = fmt.Errorf("title cannot be empty")
			return nil, txErr
		}

		var createdNote models.Note

		err := tx.QueryRow(ctx, "INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id, created_at", note.Title, note.Content).Scan(&createdNote.ID, &createdNote.CreatedAt)
		if err != nil {
			txErr = err
			return nil, txErr
		}
		createdNote.Title = note.Title
		createdNote.Content = note.Content
		created = append(created, createdNote)
	}

	if err := tx.Commit(ctx); err != nil {
		txErr = err
		return nil, txErr
	}
	return created, nil
}
func (s *PostgresStore) CreateManyForUser(ctx context.Context, userID string, notes []models.Note) ([]models.Note, error) {
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	var txErr error

	defer func() {
		if txErr != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	created := make([]models.Note, 0, len(notes))

	for _, note := range notes {

		if note.Title == "" {
			txErr = fmt.Errorf("title cannot be empty")
			return nil, txErr
		}

		var createdNote models.Note

		err := tx.QueryRow(ctx, "INSERT INTO notes (title, content, user_id) VALUES ($1, $2, $3) RETURNING id, created_at, user_id", note.Title, note.Content, userID).Scan(&createdNote.ID, &createdNote.CreatedAt, &createdNote.UserID)
		if err != nil {
			txErr = err
			return nil, txErr
		}
		createdNote.Title = note.Title
		createdNote.Content = note.Content
		created = append(created, createdNote)
	}

	if err := tx.Commit(ctx); err != nil {
		txErr = err
		return nil, txErr
	}
	return created, nil
}
