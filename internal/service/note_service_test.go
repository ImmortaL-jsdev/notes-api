package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/service/mocks"
	"go.uber.org/mock/gomock"
)

func TestNoteService_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	svc := NewNoteService(mockRepo)

	expectedNotes := []models.Note{{ID: "1", Title: "test", Content: "content"}}

	mockRepo.EXPECT().GetAll(gomock.Any()).Return(expectedNotes, nil)

	got, err := svc.GetAll(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expectedNotes) {
		t.Errorf("got %+v, want %+v", got, expectedNotes)
	}
}
func TestNoteService_GetAll_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)
	svc := NewNoteService(mockRepo)

	mockRepo.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("db error"))

	got, err := svc.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func TestNoteService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	inputNotes := models.Note{Title: "Hello", Content: "World"}

	expectedNotes := models.Note{ID: "123", Title: "Hello", Content: "World", CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

	mockRepo.EXPECT().Create(gomock.Any(), inputNotes).Return(expectedNotes, nil)

	svc := NewNoteService(mockRepo)

	got, err := svc.Create(context.Background(), inputNotes)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expectedNotes) {
		t.Errorf("got %+v, want %+v", got, expectedNotes)
	}
}

func TestNoteService_Create_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	inputNotes := models.Note{Title: "Hello", Content: "World"}

	svc := NewNoteService(mockRepo)

	mockRepo.EXPECT().Create(gomock.Any(), inputNotes).Return(models.Note{}, errors.New("db error"))

	got, err := svc.Create(context.Background(), inputNotes)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if got.ID != "" {
		t.Errorf("got %v, want nil", got)
	}

}

func TestNoteService_GetByID(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	svc := NewNoteService(mockRepo)

	inputNotes := models.Note{ID: "1", Title: "test", Content: "content"}

	expectedNotes := models.Note{ID: "1", Title: "Hello", Content: "World", CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

	mockRepo.EXPECT().GetByID(gomock.Any(), inputNotes.ID).Return(expectedNotes, nil)

	got, err := svc.GetByID(context.Background(), inputNotes.ID)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expectedNotes) {
		t.Errorf("got %+v, want %+v", got, expectedNotes)
	}
}

func TestNoteService_GetByID_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	svc := NewNoteService(mockRepo)

	testID := "non-existent"

	mockRepo.EXPECT().GetByID(gomock.Any(), testID).Return(models.Note{}, &myerrors.NotFoundError{Entity: "note", ID: testID})

	got, err := svc.GetByID(context.Background(), testID)

	var notFound *myerrors.NotFoundError

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %T: %v", err, err)
	}

	if got.ID != "" {
		t.Errorf("got %v, want nil", got)
	}
}

func TestNoteService_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	svc := NewNoteService(mockRepo)

	idTest := "1"

	inputNotes := models.Note{Title: "UpdatedTitle", Content: "UpdatedContent"}

	expectedNotes := models.Note{ID: idTest, Title: "UpdatedTitle", Content: "UpdatedContent", CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}

	mockRepo.EXPECT().Update(gomock.Any(), idTest, inputNotes).Return(expectedNotes, nil)

	got, err := svc.Update(context.Background(), idTest, inputNotes)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expectedNotes) {
		t.Errorf("got %+v, want %+v", got, expectedNotes)
	}
}
func TestNoteService_Delete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	svc := NewNoteService(mockRepo)

	idTest := "1"

	mockRepo.EXPECT().Delete(gomock.Any(), idTest).Return(nil)

	err := svc.Delete(context.Background(), idTest)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
func TestNoteService_Delete_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	mockRepo := mocks.NewMockNoteRepository(ctrl)

	svc := NewNoteService(mockRepo)

	testID := "non-existent"

	mockRepo.EXPECT().Delete(gomock.Any(), testID).Return(&myerrors.NotFoundError{Entity: "note", ID: testID})

	err := svc.Delete(context.Background(), testID)

	if err == nil {
		t.Fatal("expected error got nil")
	}

	var notFound *myerrors.NotFoundError
	if !errors.As(err, &notFound) {
		t.Fatalf("expected NotFoundError, got %T: %v", err, err)
	}
}
