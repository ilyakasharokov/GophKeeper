package service

import (
	"context"
	"gophkeeper/internal/app/server/authorization"
	"gophkeeper/internal/app/server/interfaces"
	"gophkeeper/internal/common/models"
)

// NewUserService функция создания сервиса для работы с пользователями
func NewUserService(db interfaces.DBModel, accessTokenLiveTimeMinutes int,
	refreshTokenLiveTimeDays int, accessTokenSecret string,
	refreshTokenSecret string) *UserService {
	return &UserService{
		db:                         db,
		AccessTokenLiveTimeMinutes: accessTokenLiveTimeMinutes,
		RefreshTokenLiveTimeDays:   refreshTokenLiveTimeDays,
		AccessTokenSecret:          accessTokenSecret,
		RefreshTokenSecret:         refreshTokenSecret,
	}
}

// UserService структура для сервиса пользователей
type UserService struct {
	db                         interfaces.DBModel
	AccessTokenLiveTimeMinutes int
	RefreshTokenLiveTimeDays   int
	AccessTokenSecret          string
	RefreshTokenSecret         string
}

// CreateUser функия создания пользователя
func (us *UserService) CreateUser(ctx context.Context, user models.User) error {
	return us.db.CreateUser(ctx, user)
}

// AuthUser функия авторизации пользователя
func (us *UserService) AuthUser(ctx context.Context, user models.User) (*authorization.TokenInfo, error) {
	userID, err := us.db.CheckUserPassword(ctx, user)
	if err != nil {
		return nil, err
	}
	return authorization.CreateToken(userID, us.AccessTokenLiveTimeMinutes, us.RefreshTokenLiveTimeDays,
		us.AccessTokenSecret, us.RefreshTokenSecret)
}

// RefreshToken функиця по обновлению токенов пользователя
func (us *UserService) RefreshToken(_ context.Context, refreshToken string) (*authorization.TokenInfo, error) {
	return authorization.RefreshToken(refreshToken, us.AccessTokenLiveTimeMinutes, us.RefreshTokenLiveTimeDays,
		us.AccessTokenSecret, us.RefreshTokenSecret)
}
