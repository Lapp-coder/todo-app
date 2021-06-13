package service

import (
	"errors"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

const (
	envSigningKey string = "SIGNING_KEY"
	tokenTTL             = time.Hour * 12
)

type AuthService struct {
	repos repository.Authorization
}

func NewAuthService(repos repository.Authorization) *AuthService {
	return &AuthService{repos: repos}
}

func generateHashPassword(password string) string {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(passwordHash)
}

func compareHashAndPassword(passwordHash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}

func (s AuthService) CreateUser(user model.User) (int, error) {
	user.Password = generateHashPassword(user.Password)

	return s.repos.CreateUser(user)
}

func (s AuthService) GenerateToken(email, password string) (string, error) {
	user, err := s.repos.GetUser(email)
	if err != nil || !compareHashAndPassword(user.Password, password) {
		return "", errors.New("incorrect email or password")
	}

	type tokenClaims struct {
		jwt.StandardClaims
		UserId int `json:"user_id"`
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(os.Getenv(envSigningKey)))
}

func (s AuthService) ParseToken(accessToken string) (int, error) {
	type tokenClaims struct {
		jwt.StandardClaims
		UserId int `json:"user_id"`
	}

	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(os.Getenv(envSigningKey)), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, err
	}

	return claims.UserId, nil
}
