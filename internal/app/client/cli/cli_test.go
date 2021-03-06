package cli

import (
	"context"
	"gophkeeper/internal/mocks"
	"gophkeeper/pkg/models"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestCLI_Start(t *testing.T) {
	type fields struct {
		state   int
		auther  Authenticator
		storage NotesKeeper
		user    *models.User
	}
	type args struct {
		cancel context.CancelFunc
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
			cli := &CLI{
				state:   tt.fields.state,
				auther:  tt.fields.auther,
				storage: tt.fields.storage,
				user:    tt.fields.user,
			}
			if err := cli.Start(tt.args.cancel); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLI_addNote(t *testing.T) {
	type fields struct {
		state   int
		auther  Authenticator
		storage NotesKeeper
		user    *models.User
	}
	gmock := new(mocks.GRPCClientModel)
	storage := new(mocks.NotesKeeper)
	storage.On("AddNote", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	storage.On("Flush", mock.Anything).Return(nil)
	usr := &models.User{
		ID:           "sadfdasf",
		Login:        "ilya",
		PasswordHash: []byte("qeqweqwewqeqqw"),
		Token:        "asdsadas dasd sadasdas dasda",
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "ok", fields: struct {
			state   int
			auther  Authenticator
			storage NotesKeeper
			user    *models.User
		}{state: 0, auther: gmock, storage: storage, user: usr}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{
				state:   tt.fields.state,
				auther:  tt.fields.auther,
				storage: tt.fields.storage,
				user:    tt.fields.user,
			}
			if err := cli.addNote(); (err != nil) != tt.wantErr {
				t.Errorf("addNote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCLI_auth(t *testing.T) {
	type fields struct {
		state   int
		auther  Authenticator
		storage NotesKeeper
		user    *models.User
		syncer  Syncer
	}
	gmock := new(mocks.GRPCClientModel)
	storage := new(mocks.NotesKeeper)
	gmock.On("Login", "", "").Return("", nil)
	gmock.On("SyncData", mock.Anything, mock.Anything).Return([]models.Note{}, time.Now(), nil)
	storage.On("CheckFile").Return(false)
	storage.On("MergeNotes", []models.Note{}, mock.Anything).Return(nil)
	storage.On("Flush", mock.Anything).Return(nil)
	syncer := new(mocks.Syncer)
	syncer.On("Sync").Return(nil)
	usr := models.User{
		ID:    "",
		Login: "",
		Token: "",
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "ok", fields: fields{
			state:   0,
			auther:  gmock,
			storage: storage,
			user:    &usr,
			syncer:  syncer,
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{
				state:   tt.fields.state,
				auther:  tt.fields.auther,
				storage: tt.fields.storage,
				user:    tt.fields.user,
				syncer:  tt.fields.syncer,
			}
			_, err := cli.auth()
			if (err != nil) != tt.wantErr {
				t.Errorf("auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCLI_loggedIn(t *testing.T) {
	type fields struct {
		state   int
		auther  Authenticator
		storage NotesKeeper
		user    *models.User
	}
	type args struct {
		cancel context.CancelFunc
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantQuit bool
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{
				state:   tt.fields.state,
				auther:  tt.fields.auther,
				storage: tt.fields.storage,
				user:    tt.fields.user,
			}
			gotQuit, err := cli.loggedIn(tt.args.cancel)
			if (err != nil) != tt.wantErr {
				t.Errorf("loggedIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuit != tt.wantQuit {
				t.Errorf("loggedIn() gotQuit = %v, want %v", gotQuit, tt.wantQuit)
			}
		})
	}
}

func TestCLI_registration(t *testing.T) {
	type fields struct {
		state   int
		auther  Authenticator
		storage NotesKeeper
		user    *models.User
	}
	tests := []struct {
		name    string
		fields  fields
		want    models.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{
				state:   tt.fields.state,
				auther:  tt.fields.auther,
				storage: tt.fields.storage,
				user:    tt.fields.user,
			}
			got, err := cli.registration()
			if (err != nil) != tt.wantErr {
				t.Errorf("registration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("registration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCLI_start(t *testing.T) {
	type fields struct {
		state   int
		auther  Authenticator
		storage NotesKeeper
		user    *models.User
	}
	type args struct {
		cancel context.CancelFunc
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantQuit bool
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &CLI{
				state:   tt.fields.state,
				auther:  tt.fields.auther,
				storage: tt.fields.storage,
				user:    tt.fields.user,
			}
			gotQuit, err := cli.start(tt.args.cancel)
			if (err != nil) != tt.wantErr {
				t.Errorf("start() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuit != tt.wantQuit {
				t.Errorf("start() gotQuit = %v, want %v", gotQuit, tt.wantQuit)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		auther Authenticator
		s      NotesKeeper
		syncer Syncer
	}
	tests := []struct {
		name string
		args args
		want *CLI
	}{
		{name: "ok", args: args{
			auther: nil,
			s:      nil,
		}, want: &CLI{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.auther, tt.args.s, tt.args.syncer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readAuth(t *testing.T) {
	tests := []struct {
		name      string
		wantLogin string
		wantPwd   string
	}{
		{name: "ok", wantLogin: "", wantPwd: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogin, gotPwd := readAuth()
			if gotLogin != tt.wantLogin {
				t.Errorf("readAuth() gotLogin = %v, want %v", gotLogin, tt.wantLogin)
			}
			if gotPwd != tt.wantPwd {
				t.Errorf("readAuth() gotPwd = %v, want %v", gotPwd, tt.wantPwd)
			}
		})
	}
}

func Test_readString(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "ok",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readString(); got != tt.want {
				t.Errorf("readString() = %v, want %v", got, tt.want)
			}
		})
	}
}
