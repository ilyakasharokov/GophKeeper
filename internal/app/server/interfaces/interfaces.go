package interfaces

import (
	"context"
	"gophkeeper/internal/common/models"
	"time"
)

type DBModel interface {
	CreateUser(ctx context.Context, user models.User) error
	CheckUserPassword(ctx context.Context, user models.User) (string, error)
	AddNote(ctx context.Context, userID string, note models.Note) (string, error)
	GetNotes(ctx context.Context, userID string) ([]models.Note, error)
	GetUpdates(ctx context.Context, userID string, after time.Time) ([]models.Note, error)
	UpdateNote(ctx context.Context, userID string, note models.Note) error
	DeleteNote(ctx context.Context, userID string, id string) error
}