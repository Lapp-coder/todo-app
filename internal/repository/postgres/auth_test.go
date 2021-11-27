package postgres

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthPostgres_CreateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewAuthRepository(db)

	type args struct {
		user model.User
	}

	type mockBehavior func(input args)

	testCases := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		expectedId   int
		wantErr      bool
	}{
		{
			name: "OK",
			input: args{
				model.User{Name: "user", Email: "user@gmail.com", Password: "userPassword"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)
				query := fmt.Sprintf("INSERT INTO %s (.+) VALUES (.+) RETURNING id", usersTable)
				mock.ExpectQuery(query).WithArgs(input.user.Name, input.user.Email, input.user.Password).WillReturnRows(rows)
			},
			expectedId: 1,
			wantErr:    false,
		},
		{
			name: "Empty fields",
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"})
				query := fmt.Sprintf("INSERT INTO %s (.+) VALUES (.+) RETURNING id", usersTable)
				mock.ExpectQuery(query).WithArgs(input.user.Name, input.user.Email, input.user.Password).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.CreateUser(tc.input.user)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedId, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthPostgres_GetUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewAuthRepository(db)

	type args struct {
		email string
	}

	type mockBehavior func(input args)

	testCases := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		expectedUser model.User
		wantErr      bool
	}{
		{
			name:  "OK",
			input: args{email: "user@gmail.com"},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "name", "email", "password_hash"}).
					AddRow(1, "user", "user@gmail.com", "user")
				query := fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)
				mock.ExpectQuery(query).WithArgs(input.email).WillReturnRows(rows)
			},
			expectedUser: model.User{ID: 1, Name: "user", Email: "user@gmail.com", Password: "user"},
			wantErr:      false,
		},
		{
			name: "Empty field",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s WHERE (.+)", usersTable)
				mock.ExpectQuery(query).WithArgs("")
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetUser(tc.input.email)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedUser, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
