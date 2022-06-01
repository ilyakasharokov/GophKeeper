package grpcserver

import (
	"context"
	"gophkeeper/internal/app/server/authorization"
	"gophkeeper/internal/app/server/interfaces"
	"gophkeeper/internal/mocks"
	proto "gophkeeper/pkg/grpc/proto"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	usm := new(mocks.UserServiceModel)
	ssm := new(mocks.SyncServiceModel)
	type args struct {
		u  interfaces.UserServiceModel
		ss interfaces.SyncServiceModel
	}
	tests := []struct {
		name string
		args args
		want *Server
	}{
		{name: "ok", args: args{
			u:  usm,
			ss: ssm,
		}, want: &Server{
			UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
			userService:                   usm,
			syncService:                   ssm,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.u, tt.args.ss); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_CheckConn(t *testing.T) {
	type fields struct {
		UnimplementedGophKeeperServer proto.UnimplementedGophKeeperServer
		userService                   interfaces.UserServiceModel
		syncService                   interfaces.SyncServiceModel
	}
	type args struct {
		in0 context.Context
		req *proto.CCRequest
	}
	usm := new(mocks.UserServiceModel)
	ssm := new(mocks.SyncServiceModel)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRsp *proto.CCResponse
		wantErr bool
	}{
		{name: "ok", fields: struct {
			UnimplementedGophKeeperServer proto.UnimplementedGophKeeperServer
			userService                   interfaces.UserServiceModel
			syncService                   interfaces.SyncServiceModel
		}{userService: usm, syncService: ssm}, args: args{
			in0: context.Background(),
			req: &proto.CCRequest{},
		}, wantRsp: &proto.CCResponse{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				userService:                   tt.fields.userService,
				syncService:                   tt.fields.syncService,
			}
			gotRsp, err := s.CheckConn(tt.args.in0, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckConn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRsp, tt.wantRsp) {
				t.Errorf("CheckConn() gotRsp = %v, want %v", gotRsp, tt.wantRsp)
			}
		})
	}
}

func TestServer_Login(t *testing.T) {
	usm := new(mocks.UserServiceModel)
	usm.On("AuthUser", mock.Anything, mock.Anything).Return(&authorization.TokenInfo{}, nil)
	ssm := new(mocks.SyncServiceModel)

	type fields struct {
		UnimplementedGophKeeperServer proto.UnimplementedGophKeeperServer
		userService                   interfaces.UserServiceModel
		syncService                   interfaces.SyncServiceModel
	}
	type args struct {
		ctx context.Context
		req *proto.LoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRsp *proto.LoginResponse
		wantErr bool
	}{
		{name: "ok", fields: fields{
			userService: usm,
			syncService: ssm,
		}, args: args{
			ctx: context.Background(),
			req: &proto.LoginRequest{
				Login:    "test",
				Password: "test",
			},
		}, wantRsp: &proto.LoginResponse{Status: "ok"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				userService:                   tt.fields.userService,
				syncService:                   tt.fields.syncService,
			}
			gotRsp, err := s.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRsp, tt.wantRsp) {
				t.Errorf("Login() gotRsp = %v, want %v", gotRsp, tt.wantRsp)
			}
		})
	}
}

func TestServer_RefreshToken(t *testing.T) {
	usm := new(mocks.UserServiceModel)
	ssm := new(mocks.SyncServiceModel)
	usm.On("RefreshToken", context.Background(), "refreshtoken").Return(&authorization.TokenInfo{}, nil)
	type fields struct {
		UnimplementedGophKeeperServer proto.UnimplementedGophKeeperServer
		userService                   interfaces.UserServiceModel
		syncService                   interfaces.SyncServiceModel
	}
	type args struct {
		ctx context.Context
		req *proto.RefreshTokenRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *proto.LoginResponse
		wantErr bool
	}{
		{name: "ok", fields: fields{
			UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
			userService:                   usm,
			syncService:                   ssm,
		}, args: args{context.Background(), &proto.RefreshTokenRequest{RefreshToken: "refreshtoken"}}, want: &proto.LoginResponse{Status: "ok"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				userService:                   tt.fields.userService,
				syncService:                   tt.fields.syncService,
			}
			got, err := s.RefreshToken(tt.args.ctx, tt.args.req)
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

func TestServer_Register(t *testing.T) {
	usm := new(mocks.UserServiceModel)
	usm.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	usm.On("AuthUser", mock.Anything, mock.Anything).Return(&authorization.TokenInfo{}, nil)
	ssm := new(mocks.SyncServiceModel)
	type fields struct {
		UnimplementedGophKeeperServer proto.UnimplementedGophKeeperServer
		userService                   interfaces.UserServiceModel
		syncService                   interfaces.SyncServiceModel
	}
	type args struct {
		ctx context.Context
		req *proto.RegisterRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRsp *proto.RegisterResponse
		wantErr bool
	}{
		{name: "ok", fields: fields{
			UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
			userService:                   usm,
			syncService:                   ssm,
		}, args: args{context.Background(), &proto.RegisterRequest{}}, wantRsp: &proto.RegisterResponse{Status: "created"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				userService:                   tt.fields.userService,
				syncService:                   tt.fields.syncService,
			}
			gotRsp, err := s.Register(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRsp, tt.wantRsp) {
				t.Errorf("Register() gotRsp = %v, want %v", gotRsp, tt.wantRsp)
			}
		})
	}
}

func TestServer_SyncData(t *testing.T) {
	usm := new(mocks.UserServiceModel)
	usm.On("CreateUser", mock.Anything, mock.Anything).Return(nil)
	usm.On("AuthUser", mock.Anything, mock.Anything).Return(&authorization.TokenInfo{}, nil)
	ssm := new(mocks.SyncServiceModel)
	type fields struct {
		UnimplementedGophKeeperServer proto.UnimplementedGophKeeperServer
		userService                   interfaces.UserServiceModel
		syncService                   interfaces.SyncServiceModel
	}
	type args struct {
		ctx context.Context
		req *proto.SyncDataRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRsp *proto.SyncDataResponse
		wantErr bool
	}{
		{name: "ok", fields: fields{
			UnimplementedGophKeeperServer: proto.UnimplementedGophKeeperServer{},
			userService:                   usm,
			syncService:                   ssm,
		}, args: args{context.Background(), &proto.SyncDataRequest{}}, wantRsp: &proto.SyncDataResponse{Status: "unauthorized"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedGophKeeperServer: tt.fields.UnimplementedGophKeeperServer,
				userService:                   tt.fields.userService,
				syncService:                   tt.fields.syncService,
			}
			gotRsp, err := s.SyncData(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRsp, tt.wantRsp) {
				t.Errorf("SyncData() gotRsp = %v, want %v", gotRsp, tt.wantRsp)
			}
		})
	}
}

func Test_getUserFromContext(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUserFromContext(tt.args.ctx); got != tt.want {
				t.Errorf("getUserFromContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
