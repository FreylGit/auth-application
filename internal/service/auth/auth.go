package auth

import (
	"auth-application/internal/domain/models"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	userSaver     UserSaver
	userProvider  UserProvider
	tokenProvider TokenProvider
	tokenService  TokenService
}

var (
	ErrorExpiredToken     = errors.New("token expired")
	ErrorEmailAndPassword = errors.New("failed input email or password")
)

type UserSaver interface {
	SaveUser(ctx context.Context, newUser models.User) (userId int64, err error)
}
type TokenProvider interface {
	Refresh(ctx context.Context, userId int64, token string) (rToken models.RefreshToken, err error)
	SaveRefresh(ctx context.Context, userId int64, token string) error
	UpdateRefresh(ctx context.Context, userId int64, newToken string, prevToken string) error
}
type UserProvider interface {
	User(ctx context.Context, email string) (*models.User, error)
	UserById(ctx context.Context, id int64) (*models.User, error)
}
type TokenService interface {
	CreateRefresh(userId int64) string
	NewAccess(user models.User) string
}

func New(
	userSaver UserSaver,
	userProvider UserProvider,
	tokenProvider TokenProvider,
	tokenService TokenService) *Auth {
	return &Auth{
		userSaver:     userSaver,
		userProvider:  userProvider,
		tokenProvider: tokenProvider,
		tokenService:  tokenService,
	}
}

func (a *Auth) RegistrationNewUser(ctx context.Context, newUser models.NewUser) (*models.UserResponse, error) {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	user := models.User{
		Email:        newUser.Email,
		Name:         newUser.Name,
		PasswordHash: passwordHash,
	}
	userId, err := a.userSaver.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = userId

	atoken := a.tokenService.NewAccess(user)
	rtoken := a.tokenService.CreateRefresh(userId)
	err = a.tokenProvider.SaveRefresh(ctx, userId, rtoken)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{UserId: userId, RefreshToken: rtoken, AccessToken: atoken}, nil
}

func (a *Auth) RefreshToken(ctx context.Context, userId int64, refreshToken string) (id int64, rToken string, atoken string, err error) {
	findRefresh, err := a.tokenProvider.Refresh(ctx, userId, refreshToken)
	if err != nil {
		return 0, "", "", err
	}
	if findRefresh.ExpDate.Before(time.Now()) {
		return 0, "", "", ErrorExpiredToken
	}
	newRefresh := a.tokenService.CreateRefresh(userId)
	user, err := a.userProvider.UserById(ctx, userId)
	if err != nil {
		return 0, "", "", err
	}
	err = a.tokenProvider.UpdateRefresh(ctx, userId, newRefresh, findRefresh.Token)
	if err != nil {
		return 0, "", "", err
	}
	newAccess := a.tokenService.NewAccess(*user)

	return userId, newRefresh, newAccess, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string) (*models.UserResponse, error) {
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return nil, fmt.Errorf("%w", ErrorEmailAndPassword)
	}
	atoken := a.tokenService.NewAccess(*user)
	rtoken := a.tokenService.CreateRefresh(user.Id)
	err = a.tokenProvider.SaveRefresh(ctx, user.Id, rtoken)
	if err != nil {
		return nil, err
	}

	return &models.UserResponse{UserId: user.Id, RefreshToken: rtoken, AccessToken: atoken}, nil
}
