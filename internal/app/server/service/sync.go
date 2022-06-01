package service

import (
	"context"
	"gophkeeper/internal/app/server/interfaces"
	"gophkeeper/pkg/models"
	"time"

	"github.com/rs/zerolog/log"
)

type SyncService struct {
	db interfaces.DBModel
}

// NewSyncService функция создания сервиса для работы с пользователями
func NewSyncService(db interfaces.DBModel) *SyncService {
	return &SyncService{
		db: db,
	}
}

func (s *SyncService) Sync(ctx context.Context, userID string, notes []models.Note, lastSyncDate time.Time) ([]models.Note, time.Time, error) {
	mp := make(map[string]models.Note)
	for _, n := range notes {
		if n.ID == "" {
			id, err := s.db.AddNote(ctx, userID, n)
			if err != nil {
				log.Err(err).Msg("insert new note error")
			}
			if id != "" {
				n.ID = id
				mp[n.ID] = n
			}
		} else {
			if n.Deleted {
				err := s.db.DeleteNote(ctx, userID, n.ID)
				if err != nil {
					log.Err(err).Str("id", n.ID).Msg("delete note error")
				}
			} else {
				err := s.db.UpdateNote(ctx, userID, n)
				if err != nil {
					log.Err(err).Str("id", n.ID).Msg("delete note error")
				}
			}
		}
	}
	updates, err := s.db.GetNotesAfter(ctx, userID, lastSyncDate)
	if err != nil {
		log.Err(err).Msg("get updates error")
	}
	for _, u := range updates {
		if _, ok := mp[u.ID]; !ok {
			mp[u.ID] = u
		}
	}
	ret := make([]models.Note, 0)
	for _, m := range mp {
		ret = append(ret, m)
	}
	return ret, time.Now(), nil
}
