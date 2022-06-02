// Package grpc implement gRPC client
package grpcclient

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/credentials"
	proto "gophkeeper/pkg/grpc/proto"
	"gophkeeper/pkg/models"
	"gophkeeper/pkg/service"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Client struct {
	cc           *grpc.ClientConn
	service      proto.GophKeeperClient
	token        string
	refreshToken string
	tm           time.Duration
}

// New returns a new grpc client
func New(grpcAddr string, tm time.Duration, crd credentials.TransportCredentials) (*Client, error) {
	transportOption := grpc.WithTransportCredentials(crd)
	cc, err := grpc.Dial(grpcAddr, transportOption)
	if err != nil {
		log.Err(err).Msg("cannot dial server")
		return nil, err
	}
	srv := proto.NewGophKeeperClient(cc)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &proto.CCRequest{}
	_, err = srv.CheckConn(ctx, req)
	if err != nil {
		log.Err(err).Msg("cannot connect to server")
		return nil, err
	}
	return &Client{cc, srv, "", "", tm}, nil
}

// Login user and returns the access token
func (client *Client) Login(login string, pwd string) (token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.tm*time.Second)
	defer cancel()
	req := &proto.LoginRequest{
		Login:    login,
		Password: pwd,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}
	client.token = res.GetAccessToken()
	client.refreshToken = res.GetRefreshToken()
	return client.token, nil
}

// Registration create user and returns the access token
func (client *Client) Registration(login string, pwd string) (status string, token string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.tm*time.Second)
	defer cancel()

	req := &proto.RegisterRequest{
		Login:    login,
		Password: pwd,
	}

	res, err := client.service.Register(ctx, req)
	if err != nil {
		return "", "", err
	}
	client.token = res.GetAccessToken()
	client.refreshToken = res.GetRefreshToken()
	return res.GetStatus(), client.token, nil
}

func (client *Client) SyncData(notes []models.Note, lastSync time.Time) ([]models.Note, time.Time, error) {
	ctx, cancel := context.WithTimeout(context.Background(), client.tm*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", client.token))
	var syncDate time.Time
	gdata := service.NotesToProto(notes)
	req := &proto.SyncDataRequest{
		Notes:    gdata,
		LastSync: timestamppb.New(lastSync),
	}
	res, err := client.service.SyncData(ctx, req)
	if res != nil && res.Status == "unauthorized" {
		accessToken, refreshToken, errr := client.RefreshToken(ctx, client.refreshToken)
		if errr != nil {
			return nil, syncDate, errr
		}
		client.token = accessToken
		client.refreshToken = refreshToken
		ctx, cancel := context.WithTimeout(context.Background(), client.tm*time.Second)
		defer cancel()
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", client.token))
		res, err = client.service.SyncData(ctx, req)
	}
	if err != nil {
		return nil, syncDate, err
	}
	syncDate = res.GetLastSync().AsTime()
	gdata = res.GetNotes()
	notes = service.ProtoNotesToModels(gdata)
	return notes, syncDate, nil
}

// RefreshToken функция обновления токена пользователя
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	message := proto.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}
	response, err := c.service.RefreshToken(ctx, &message)
	if err != nil {
		return "", "", err
	}
	if response.Status == "ok" {
		return response.AccessToken, response.RefreshToken, nil
	}
	return "", "", errors.New(response.Status)
}

func (client *Client) Close() {
	client.cc.Close()
}
