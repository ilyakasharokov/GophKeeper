package interfaces

import (
	"context"
	"gophkeeper/internal/app/server/authorization"
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

type UserServiceModel interface {
	RefreshToken(ctx context.Context, token string) (*authorization.TokenInfo, error)
	CreateUser(ctx context.Context, user models.User) error
	AuthUser(ctx context.Context, user models.User) (*authorization.TokenInfo, error)
}

type SyncServiceModel interface {
	Sync(ctx context.Context, userID string, notes []models.Note, lastSyncDate time.Time) ([]models.Note, time.Time, error)
}
