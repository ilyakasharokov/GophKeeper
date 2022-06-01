package service

import (
	"context"
	"gophkeeper/internal/app/server/authorization"
	"gophkeeper/internal/app/server/interfaces"
	"gophkeeper/internal/mocks"
	"gophkeeper/pkg/models"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestNewSyncService(t *testing.T) {
	type args struct {
		db interfaces.DBModel
	}
	pgdb := new(mocks.DBModel)
	tests := []struct {
		name string
		args args
		want *SyncService
	}{
		{name: "ok", args: args{db: pgdb}, want: &SyncService{db: pgdb}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSyncService(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSyncService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUserService(t *testing.T) {
	type args struct {
		db                         interfaces.DBModel
		accessTokenLiveTimeMinutes int
		refreshTokenLiveTimeDays   int
		accessTokenSecret          string
		refreshTokenSecret         string
	}
	db := new(mocks.DBModel)
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		{name: "ok", args: args{
			db:                         db,
			accessTokenLiveTimeMinutes: 1,
			refreshTokenLiveTimeDays:   1,
			accessTokenSecret:          "ats",
			refreshTokenSecret:         "rts",
		}, want: &UserService{
			db:                         db,
			AccessTokenLiveTimeMinutes: 1,
			RefreshTokenLiveTimeDays:   1,
			AccessTokenSecret:          "ats",
			RefreshTokenSecret:         "rts",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.db, tt.args.accessTokenLiveTimeMinutes, tt.args.refreshTokenLiveTimeDays, tt.args.accessTokenSecret, tt.args.refreshTokenSecret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncService_Sync(t *testing.T) {
	type fields struct {
		db interfaces.DBModel
	}
	type args struct {
		ctx          context.Context
		userID       string
		notes        []models.Note
		lastSyncDate time.Time
	}
	// db := new(mocks.DBModel)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []models.Note
		want1   time.Time
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncService{
				db: tt.fields.db,
			}
			got, got1, err := s.Sync(tt.args.ctx, tt.args.userID, tt.args.notes, tt.args.lastSyncDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sync() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Sync() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestUserService_AuthUser(t *testing.T) {
	type fields struct {
		db                         interfaces.DBModel
		AccessTokenLiveTimeMinutes int
		RefreshTokenLiveTimeDays   int
		AccessTokenSecret          string
		RefreshTokenSecret         string
	}
	type args struct {
		ctx  context.Context
		user models.User
	}
	db := new(mocks.DBModel)
	db.On("CheckUserPassword", mock.Anything, mock.Anything).Return("true", nil)
	token, _ := authorization.CreateToken("", 0, 0,
		"", "")
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *authorization.TokenInfo
		wantErr bool
	}{
		{name: "ok", fields: fields{
			db: db,
		}, args: args{context.Background(), models.User{}}, want: token, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db:                         tt.fields.db,
				AccessTokenLiveTimeMinutes: tt.fields.AccessTokenLiveTimeMinutes,
				RefreshTokenLiveTimeDays:   tt.fields.RefreshTokenLiveTimeDays,
				AccessTokenSecret:          tt.fields.AccessTokenSecret,
				RefreshTokenSecret:         tt.fields.RefreshTokenSecret,
			}
			_, err := us.AuthUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*			if !reflect.DeepEqual(got, tt.want) {
						t.Errorf("AuthUser() got = %v, want %v", got, tt.want)
					}*/
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	type fields struct {
		db                         interfaces.DBModel
		AccessTokenLiveTimeMinutes int
		RefreshTokenLiveTimeDays   int
		AccessTokenSecret          string
		RefreshTokenSecret         string
	}
	type args struct {
		ctx  context.Context
		user models.User
	}
	db := new(mocks.DBModel)
	db.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "ok", fields: fields{
			db: db,
		}, args: args{context.Background(), models.User{}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db:                         tt.fields.db,
				AccessTokenLiveTimeMinutes: tt.fields.AccessTokenLiveTimeMinutes,
				RefreshTokenLiveTimeDays:   tt.fields.RefreshTokenLiveTimeDays,
				AccessTokenSecret:          tt.fields.AccessTokenSecret,
				RefreshTokenSecret:         tt.fields.RefreshTokenSecret,
			}
			if err := us.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_RefreshToken(t *testing.T) {
	type fields struct {
		db                         interfaces.DBModel
		AccessTokenLiveTimeMinutes int
		RefreshTokenLiveTimeDays   int
		AccessTokenSecret          string
		RefreshTokenSecret         string
	}
	type args struct {
		in0          context.Context
		refreshToken string
	}
	db := new(mocks.DBModel)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *authorization.TokenInfo
		wantErr bool
	}{
		{name: "ok", fields: struct {
			db                         interfaces.DBModel
			AccessTokenLiveTimeMinutes int
			RefreshTokenLiveTimeDays   int
			AccessTokenSecret          string
			RefreshTokenSecret         string
		}{db: db, AccessTokenLiveTimeMinutes: 1, RefreshTokenLiveTimeDays: 1, AccessTokenSecret: "1", RefreshTokenSecret: "2"}, args: struct {
			in0          context.Context
			refreshToken string
		}{in0: context.Background(), refreshToken: "1"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UserService{
				db:                         tt.fields.db,
				AccessTokenLiveTimeMinutes: tt.fields.AccessTokenLiveTimeMinutes,
				RefreshTokenLiveTimeDays:   tt.fields.RefreshTokenLiveTimeDays,
				AccessTokenSecret:          tt.fields.AccessTokenSecret,
				RefreshTokenSecret:         tt.fields.RefreshTokenSecret,
			}
			got, err := us.RefreshToken(tt.args.in0, tt.args.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
