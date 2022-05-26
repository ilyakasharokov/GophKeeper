package grpcserver

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gophkeeper/internal/app/server/service"
	customerrors "gophkeeper/internal/common/errors"
	"gophkeeper/internal/common/models"
	service2 "gophkeeper/internal/common/service"
	proto "gophkeeper/pkg/grpc/proto"
)

// Server implement main methods for gRPC
type Server struct {
	proto.UnimplementedGophKeeperServer
	userService *service.UserService
	syncService *service.SyncService
}

// New instance for gRPC server
func New(u *service.UserService, ss *service.SyncService) *Server {
	return &Server{
		userService: u,
		syncService: ss,
	}
}

func (s *Server) Login(ctx context.Context, req *proto.LoginRequest) (rsp *proto.LoginResponse, err error) {
	rsp = new(proto.LoginResponse)
	user := models.User{
		Login:    req.Login,
		Password: req.Password,
	}
	tokens, err := s.userService.AuthUser(ctx, user)
	if err != nil {
		return &proto.LoginResponse{
			Status: customerrors.ParseError(err),
		}, err
	}
	return &proto.LoginResponse{
		Status:       "ok",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// RefreshToken функция для обновления токена пользователя
func (s *Server) RefreshToken(ctx context.Context, req *proto.RefreshTokenRequest) (*proto.LoginResponse, error) {
	tokens, err := s.userService.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return &proto.LoginResponse{
			Status: customerrors.ParseError(err),
		}, err
	}
	return &proto.LoginResponse{
		Status:       "ok",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *Server) Register(ctx context.Context, req *proto.RegisterRequest) (rsp *proto.RegisterResponse, err error) {
	user := models.User{
		Login:    req.Login,
		Password: req.Password,
	}
	err = s.userService.CreateUser(ctx, user)
	if err != nil {
		return &proto.RegisterResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	tokens, err := s.userService.AuthUser(ctx, user)
	if err != nil {
		return &proto.RegisterResponse{
			Status: customerrors.ParseError(err),
		}, nil
	}
	return &proto.RegisterResponse{
		Status:       "created",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *Server) CheckConn(_ context.Context, req *proto.CCRequest) (rsp *proto.CCResponse, err error) {
	rsp = new(proto.CCResponse)
	return rsp, nil
}

func (s *Server) SyncData(ctx context.Context, req *proto.SyncDataRequest) (rsp *proto.SyncDataResponse, err error) {
	userID := getUserFromContext(ctx)
	rsp = new(proto.SyncDataResponse)
	if userID == "" {
		rsp.Status = "unauthorized"
		return rsp, nil
	}
	notes := service2.ProtoNotesToModels(req.GetNotes())
	retNotes, retLastSyncDate, err := s.syncService.Sync(ctx, userID, notes, req.LastSync.AsTime())
	rsp.Notes = service2.NotesToProto(retNotes)
	rsp.LastSync = timestamppb.New(retLastSyncDate)
	return rsp, nil
}

/*func getUserID(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	userID := "default"
	var uid []string = nil
	if ok {
		uid = md.Get("user_id")
	}
	if len(uid) > 0 {
		userID = uid[0]
	} else {
		userID = uuid.New().String()
	}
	return userID
}
*/

func getUserFromContext(ctx context.Context) string {
	userID := ctx.Value("userID")
	if userID != nil {
		if str, ok := userID.(string); ok {
			return str
		} else {
			return ""
		}
	}
	return ""
}
