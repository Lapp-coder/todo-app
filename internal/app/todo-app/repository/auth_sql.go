package repository

import (
	"errors"
	"fmt"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/jmoiron/sqlx"
)

type AuthSQL struct {
	db *sqlx.DB
}

func NewAuthSQL(db *sqlx.DB) *AuthSQL {
	return &AuthSQL{db: db}
}

func (r AuthSQL) CreateUser(user model.User) (int, error) {
	err := r.db.QueryRow(fmt.Sprintf(
		"INSERT INTO %s (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id", usersTable),
		user.Name, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return 0, errors.New("failed to create account")
	}

	return user.Id, nil
}

func (r AuthSQL) GetByEmail(email string) (model.User, error) {
	var user model.User

	err := r.db.QueryRow(fmt.Sprintf(
		"SELECT id, name, email, password_hash FROM %s WHERE email = $1", usersTable),
		email).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	return user, err
}
