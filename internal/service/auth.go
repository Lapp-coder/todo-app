package service

import (
	"time"

	"github.com/Lapp-coder/todo-app/internal/config"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/repository"
	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	repos repository.Authorization
	cfg   config.Service
}

func NewAuthService(repos repository.Authorization, cfg config.Service) *AuthService {
	return &AuthService{repos: repos, cfg: cfg}
}

func (s AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generatePasswordHash(user.Password, s.cfg.Salt)

	return s.repos.CreateUser(user)
}

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

func (s AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repos.GetUser(email)
	if err != nil || !compareHashAndPassword(user.Password, password, s.cfg.Salt) {
		return "", ErrIncorrectEmailOrPassword
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Second * time.Duration(s.cfg.TokenTTL)).Unix(),
		},
		user.ID,
	})

	return token.SignedString([]byte(s.cfg.SigningKey))
}

func (s AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSigningMethod
		}

		return []byte(s.cfg.SigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, err
	}

	return claims.UserID, nil
}
