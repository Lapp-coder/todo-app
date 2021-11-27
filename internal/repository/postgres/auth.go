package postgres

import (
	"fmt"

	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/jmoiron/sqlx"
)

type AuthRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(user model.User) (int, error) {
	if err := r.db.QueryRow(fmt.Sprintf(
		"INSERT INTO %s (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable),
		user.Name, user.Email, user.Password).Scan(&user.ID); err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (r *AuthRepository) GetUser(email string) (model.User, error) {
	var user model.User
	if err := r.db.QueryRow(fmt.Sprintf(
		"SELECT id, name, email, password_hash FROM %s WHERE email = $1", usersTable),
		email).Scan(&user.ID, &user.Name, &user.Email, &user.Password); err != nil {
		return model.User{}, err
	}

	return user, nil
}
