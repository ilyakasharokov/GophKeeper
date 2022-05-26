package authorization

import (
	"context"
	"github.com/golang-jwt/jwt"
	"reflect"
	"testing"
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
		want    *TokenDetails
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateToken(tt.args.userID, tt.args.accessTokenLiveTimeMinutes, tt.args.refreshTokenLiveTimeDays, tt.args.accessTokenSecret, tt.args.refreshTokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractToken(t *testing.T) {
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
		want    *TokenDetails
		wantErr bool
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
