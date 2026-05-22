package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ImmortaL-jsdev/notes-api/internal/auth"
	myerrors "github.com/ImmortaL-jsdev/notes-api/internal/errors"
	"github.com/ImmortaL-jsdev/notes-api/internal/models"
	"github.com/ImmortaL-jsdev/notes-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  *repository.UserStore
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserStore, secret []byte) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: secret,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (models.User, error) {

	if email == "" || password == "" {
		return models.User{}, &myerrors.ValidationError{Message: "email and password are required"}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	created, err := s.userRepo.CreateUser(ctx, email, string(hashedPassword))

	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	created.PasswordHash = ""

	return created, err
}

func (s *AuthService) Login(ctx context.Context, email, password string) (models.TokenPair, error) {

	if email == "" || password == "" {
		return models.TokenPair{}, &myerrors.ValidationError{Message: "email and password are required"}
	}

	user, err := s.userRepo.GetUserByEmail(ctx, email)

	if err != nil {
		return models.TokenPair{}, &myerrors.ValidationError{Message: "invalid credentials"}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return models.TokenPair{}, &myerrors.ValidationError{Message: "invalid credentials"}
	}

	accessToken, err := auth.GenerateToken(user.ID, 15*time.Minute, s.jwtSecret)

	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed access token: %w", err)
	}

	refreshToken, err := auth.GenerateToken(user.ID, 7*24*time.Hour, s.jwtSecret)

	if err != nil {
		return models.TokenPair{}, fmt.Errorf("failed refresh token: %w", err)
	}

	return models.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
