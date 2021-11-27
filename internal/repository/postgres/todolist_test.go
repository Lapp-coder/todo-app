package postgres

import (
	"fmt"
	"testing"

	"github.com/Lapp-coder/todo-app/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoListPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoListRepository(db)

	type args struct {
		userID int
		list   model.TodoList
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name           string
		input          args
		mockBehavior   mockBehavior
		expectedListID int
		wantErr        bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				userID: 1,
				list:   model.TodoList{UserID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)

				query := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID, input.list.Title, input.list.Description, input.list.CompletionDate).WillReturnRows(rows)
			},
			expectedListID: 1,
			wantErr:        false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				userID: 1,
				list:   model.TodoList{Title: "test", CompletionDate: "2021-11-21 00:00:00"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)

				query := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID, input.list.Title, input.list.CompletionDate).WillReturnRows(rows)
			},
			expectedListID: 1,
			wantErr:        false,
		},
		{
			name: "OK_WithoutCompletionDate",
			input: args{
				userID: 1,
				list:   model.TodoList{Title: "test", Description: "testing"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)

				query := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID, input.list.Title, input.list.Description).WillReturnRows(rows)
			},
			expectedListID: 1,
			wantErr:        false,
		},
		{
			name: "Empty fields",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID)
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.Create(tc.input.userID, tc.input.list)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedListID, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_GetAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoListRepository(db)

	type args struct {
		userID int
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name          string
		input         args
		mockBehavior  mockBehavior
		expectedLists []model.TodoList
		wantErr       bool
	}{
		{
			name:  "OK",
			input: args{userID: 1},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "user_id", "title", "description", "completion_date"}).
					AddRow(1, 1, "test1", "testing1", "2021-11-12 00:00:00").
					AddRow(2, 1, "test2", "testing2", "2021-11-12 00:00:00").
					AddRow(3, 1, "test3", "testing3", "2021-11-12 00:00:00")

				query := fmt.Sprintf("SELECT (.+) FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID).WillReturnRows(rows)
			},
			expectedLists: []model.TodoList{
				{ID: 1, UserID: 1, Title: "test1", Description: "testing1", CompletionDate: "2021-11-12 00:00:00"},
				{ID: 2, UserID: 1, Title: "test2", Description: "testing2", CompletionDate: "2021-11-12 00:00:00"},
				{ID: 3, UserID: 1, Title: "test3", Description: "testing3", CompletionDate: "2021-11-12 00:00:00"},
			},
			wantErr: false,
		},
		{
			name: "Empty field",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectQuery(query).WithArgs(0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetAll(tc.input.userID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedLists, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_GetByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoListRepository(db)

	type args struct {
		userID int
		listID int
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		expectedList model.TodoList
		wantErr      bool
	}{
		{
			name: "OK",
			input: args{
				userID: 1,
				listID: 3,
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "user_id", "title", "description", "completion_date"}).
					AddRow(3, 1, "test", "testing", "2021-11-12 00:00:00")
				query := fmt.Sprintf("SELECT (.+) FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID, input.listID).WillReturnRows(rows)
			},
			expectedList: model.TodoList{ID: 3, UserID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-12 00:00:00"},
			wantErr:      false,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectQuery(query).WithArgs(0, 0)
			},
			wantErr: true,
		},
		{
			name:  "EmptyField_UserId",
			input: args{listID: 3},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectQuery(query).WithArgs(0, 3)
			},
			wantErr: true,
		},
		{
			name:  "EmptyField_ListId",
			input: args{userID: 1},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectQuery(query).WithArgs(1, 0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetByID(tc.input.userID, tc.input.listID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedList, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoListRepository(db)

	type args struct {
		listID int
		update model.UpdateTodoList
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				listID: 1,
				update: model.UpdateTodoList{
					Title:          test.StringPointer("test"),
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, *input.update.CompletionDate, input.listID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				listID: 1,
				update: model.UpdateTodoList{
					Title:          test.StringPointer("test"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.CompletionDate, input.listID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitle",
			input: args{
				listID: 1,
				update: model.UpdateTodoList{
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Description, *input.update.CompletionDate, input.listID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutCompletionDate",
			input: args{
				listID: 1,
				update: model.UpdateTodoList{
					Title:       test.StringPointer("test"),
					Description: test.StringPointer("testing"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, input.listID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_NoUpdateFields",
			input: args{
				listID: 1,
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(input.listID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Update(tc.input.listID, tc.input.update)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListPostgres_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoListRepository(db)

	type args struct {
		listID int
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:  "OK",
			input: args{listID: 1},
			mockBehavior: func(input args) {
				mock.ExpectBegin()

				query1 := fmt.Sprintf("DELETE FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query1).WithArgs(input.listID).WillReturnResult(sqlmock.NewResult(1, 1))

				query2 := fmt.Sprintf("DELETE FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectExec(query2).WithArgs(input.listID).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "Empty field",
			mockBehavior: func(input args) {
				mock.ExpectBegin()

				query1 := fmt.Sprintf("DELETE FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query1).WithArgs(0).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Delete(tc.input.listID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
