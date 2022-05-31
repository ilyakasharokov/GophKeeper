package authorization

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/golang-jwt/jwt"
)

func TestCreateToken(t *testing.T) {
	type args struct {
		userID                     string
		accessTokenLiveTimeMinutes int
		refreshTokenLiveTimeDays   int
		accessTokenSecret          string
		refreshTokenSecret         string
	}
	tests := []struct {
		name    string
		args    args
		want    *TokenInfo
		wantErr bool
	}{
		{name: "ok", args: struct {
			userID                     string
			accessTokenLiveTimeMinutes int
			refreshTokenLiveTimeDays   int
			accessTokenSecret          string
			refreshTokenSecret         string
		}{userID: "", accessTokenLiveTimeMinutes: 0, refreshTokenLiveTimeDays: 0, accessTokenSecret: "", refreshTokenSecret: ""}, want: &TokenInfo{}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateToken(tt.args.userID, tt.args.accessTokenLiveTimeMinutes, tt.args.refreshTokenLiveTimeDays, tt.args.accessTokenSecret, tt.args.refreshTokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctxWithToken := metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("Bearer %v", "token"))
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "ok", args: args{ctx: context.Background()}, want: ""},
		{name: "ok", args: args{ctx: ctxWithToken}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractToken(tt.args.ctx); got != tt.want {
				t.Errorf("ExtractToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRefreshToken(t *testing.T) {
	type args struct {
		refresh                    string
		accessTokenLiveTimeMinutes int
		refreshTokenLiveTimeDays   int
		accessTokenSecret          string
		refreshTokenSecret         string
	}
	tests := []struct {
		name    string
		args    args
		want    *TokenInfo
		wantErr bool
	}{
		{name: "ok", args: args{
			refresh:                    "",
			accessTokenLiveTimeMinutes: 0,
			refreshTokenLiveTimeDays:   0,
			accessTokenSecret:          "",
			refreshTokenSecret:         "",
		}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RefreshToken(tt.args.refresh, tt.args.accessTokenLiveTimeMinutes, tt.args.refreshTokenLiveTimeDays, tt.args.accessTokenSecret, tt.args.refreshTokenSecret)
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

func TestTokenValid(t *testing.T) {
	type args struct {
		ctx          context.Context
		accessSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "ok", args: struct {
			ctx          context.Context
			accessSecret string
		}{ctx: context.Background(), accessSecret: "tadfasdfasdf"}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TokenValid(tt.args.ctx, tt.args.accessSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TokenValid() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type args struct {
		ctx          context.Context
		accessSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    *jwt.Token
		wantErr bool
	}{
		{name: "ok", args: struct {
			ctx          context.Context
			accessSecret string
		}{ctx: context.Background(), accessSecret: "qeqweqw"}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyToken(tt.args.ctx, tt.args.accessSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VerifyToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
