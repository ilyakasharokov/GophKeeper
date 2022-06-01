package syncer

import (
	"gophkeeper/internal/mocks"
	"gophkeeper/pkg/models"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		s      StorageModel
		client SyncClient
	}
	sm := new(mocks.StorageModel)
	sc := new(mocks.SyncClient)
	tests := []struct {
		name string
		args args
		want *Syncer
	}{
		{name: "ok", args: args{
			s:      sm,
			client: sc,
		}, want: &Syncer{
			storage:    sm,
			syncClient: sc,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.s, tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncer_GetLastSyncDate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		date := time.Now()
		s := Syncer{
			LastSyncDate: date,
		}
		if got := s.GetLastSyncDate(); got != date {
			t.Errorf("GetLastSyncDate() = %v, want %v", got, date)
		}
	})
}

func TestSyncer_GetNonSyncNotes(t *testing.T) {
	sm := new(mocks.StorageModel)
	notes := []models.Note{{}, {}, {}}
	sm.On("GetNotes", true).Return(notes)
	tests := []struct {
		name string
		want []models.Note
	}{
		{name: "ok", want: []models.Note{}},
	}
	syncer := &Syncer{
		storage:      sm,
		LastSyncDate: time.Now(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := syncer.GetNonSyncNotes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNonSyncNotes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncer_Sync(t *testing.T) {
	type fields struct {
		storage      StorageModel
		LastSyncDate time.Time
		syncClient   SyncClient
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Syncer{
				storage:      tt.fields.storage,
				LastSyncDate: tt.fields.LastSyncDate,
				syncClient:   tt.fields.syncClient,
			}
			if err := s.Sync(); (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSyncer_mergeNotes(t *testing.T) {
	type fields struct {
		storage      StorageModel
		LastSyncDate time.Time
		syncClient   SyncClient
	}
	type args struct {
		newdata []models.Note
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Syncer{
				storage:      tt.fields.storage,
				LastSyncDate: tt.fields.LastSyncDate,
				syncClient:   tt.fields.syncClient,
			}
			if err := s.mergeNotes(tt.args.newdata); (err != nil) != tt.wantErr {
				t.Errorf("mergeNotes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
