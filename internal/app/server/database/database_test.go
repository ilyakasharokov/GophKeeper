package database

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gophkeeper/pkg/models"
	"reflect"
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	type args struct {
		conn *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want *PGDB
	}{
		{name: "ok", args: args{conn: &sqlx.DB{}}, want: &PGDB{
			conn: &sqlx.DB{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDB(tt.args.conn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGDB_AddNote(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		userID string
		note   models.Note
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			got, err := pg.AddNote(tt.args.ctx, tt.args.userID, tt.args.note)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddNote() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddNote() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGDB_CheckUserPassword(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user models.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			got, err := pg.CheckUserPassword(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckUserPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckUserPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGDB_CreateUser(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		user models.User
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
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			if err := pg.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGDB_DeleteNote(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		userID string
		id     string
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
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			if err := pg.DeleteNote(tt.args.ctx, tt.args.userID, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGDB_GetNotes(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Note
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			got, err := pg.GetNotes(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNotes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGDB_GetUpdates(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		userID string
		after  time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Note
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			got, err := pg.GetNotesAfter(tt.args.ctx, tt.args.userID, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNonSyncNotes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNonSyncNotes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGDB_UpdateNote(t *testing.T) {
	type fields struct {
		conn *sqlx.DB
	}
	type args struct {
		ctx    context.Context
		userID string
		note   models.Note
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
			pg := &PGDB{
				conn: tt.fields.conn,
			}
			if err := pg.UpdateNote(tt.args.ctx, tt.args.userID, tt.args.note); (err != nil) != tt.wantErr {
				t.Errorf("UpdateNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
