// storage keeps data in memory and save it to disk
package syncer

import (
	"errors"
	"gophkeeper/pkg/models"
	"time"
)

type Syncer struct {
	storage      StorageModel
	LastSyncDate time.Time
	syncClient   SyncClient
}

type StorageModel interface {
	GetNotes(all bool) []models.Note
	SetNotes(notes []models.Note)
	GetLastSyncDate() time.Time
	SetLastSyncDate(date time.Time)
}

type SyncClient interface {
	SyncData(notes []models.Note, lastSync time.Time) (notesUpdate []models.Note, newLastSync time.Time, err error)
}

func New(s StorageModel, client SyncClient) *Syncer {
	return &Syncer{
		storage:    s,
		syncClient: client,
	}
}

func (s *Syncer) Sync() error {
	notesAfter := s.GetNonSyncNotes()
	remoteNotes, remoteLastSync, err := s.syncClient.SyncData(notesAfter, s.LastSyncDate)
	if err != nil {
		return errors.New("sync error (" + err.Error() + ")")
	}
	if s.LastSyncDate.After(remoteLastSync) {
		return errors.New("sync error (last sync date is incorrect")
	}
	err = s.mergeNotes(remoteNotes)
	if err != nil {
		return errors.New("merge data error (" + err.Error() + ")")
	}
	s.LastSyncDate = remoteLastSync
	s.storage.SetLastSyncDate(s.LastSyncDate)
	return nil
}

// GetNonSyncNotes get notes that were updated after last synchronization
func (s *Syncer) GetNonSyncNotes() []models.Note {
	currentNotes := s.storage.GetNotes(true)
	filtered := make([]models.Note, 0)
	for _, d := range currentNotes {
		if d.DeletedAt.After(s.LastSyncDate) || d.CreatedAt.After(s.LastSyncDate) || d.UpdatedAt.After(s.LastSyncDate) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// mergeNotes updates local notes with remote notes
func (s *Syncer) mergeNotes(newdata []models.Note) error {
	currentNotes := s.storage.GetNotes(true)
	filtered := make([]models.Note, 0)
	// удаление не синхронизированых и обновление синхронизированых
	for _, d := range currentNotes {
		if d.ID == "" || d.Deleted {
			continue
		}
		touched := false
		for _, nd := range newdata {
			if nd.ID == d.ID {
				if !nd.Deleted {
					filtered = append(filtered, nd)
				}
				touched = true
			}
		}
		if !touched {
			filtered = append(filtered, d)
		}
	}
	// добавление новых
	for _, nd := range newdata {
		if nd.CreatedAt.After(s.LastSyncDate) && !nd.Deleted {
			filtered = append(filtered, nd)
		}
	}
	s.storage.SetNotes(filtered)
	return nil
}

func (s *Syncer) UpdateLastSyncDate() {
	s.LastSyncDate = s.storage.GetLastSyncDate()
}

func (s *Syncer) GetLastSyncDate() time.Time {
	return s.LastSyncDate
}
