package repository

import (
	"fmt"
	"testing"

	"github.com/Lapp-coder/todo-app/internal/app/request"
	"github.com/Lapp-coder/todo-app/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoListSQL_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoListSQL(db)

	type args struct {
		userId int
		list   model.TodoList
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name           string
		input          args
		mockBehavior   mockBehavior
		expectedListId int
		wantErr        bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				userId: 1,
				list:   model.TodoList{Title: "test", Description: "testing", CompletionDate: "20210101:110611"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectBegin()

				query1 := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query1).WithArgs(input.list.Title, input.list.Description, input.list.CompletionDate).WillReturnRows(rows)

				query2 := fmt.Sprintf("INSERT INTO %s", usersListsTable)
				mock.ExpectExec(query2).WithArgs(input.userId, 1).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedListId: 1,
			wantErr:        false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				userId: 1,
				list:   model.TodoList{Title: "test", CompletionDate: "20210101:110611"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectBegin()

				query1 := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query1).WithArgs(input.list.Title, input.list.Description, input.list.CompletionDate).WillReturnRows(rows)

				query2 := fmt.Sprintf("INSERT INTO %s", usersListsTable)
				mock.ExpectExec(query2).WithArgs(input.userId, 1).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedListId: 1,
			wantErr:        false,
		},
		{
			name: "OK_WithoutCompletionDate",
			input: args{
				userId: 1,
				list:   model.TodoList{Title: "test"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(1)

				mock.ExpectBegin()

				query1 := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query1).WithArgs(input.list.Title, input.list.Description, input.list.CompletionDate).WillReturnRows(rows)

				query2 := fmt.Sprintf("INSERT INTO %s", usersListsTable)
				mock.ExpectExec(query2).WithArgs(input.userId, 1).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			expectedListId: 1,
			wantErr:        false,
		},
		{
			name: "Empty fields",
			mockBehavior: func(input args) {
				mock.ExpectBegin()

				query := fmt.Sprintf("INSERT INTO %s", todoListsTable)
				mock.ExpectQuery(query).WithArgs("", "", "")

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.Create(tc.input.userId, tc.input.list)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedListId, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListSQL_GetAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoListSQL(db)

	type args struct {
		userId int
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
			input: args{userId: 1},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "title", "description"}).
					AddRow(1, "test1", "testing1").
					AddRow(2, "test2", "testing2").
					AddRow(3, "test3", "testing3")

				query := fmt.Sprintf("SELECT (.+) FROM %s tl INNER JOIN %s ul ON (.+) WHERE (.+)", todoListsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(input.userId).WillReturnRows(rows)
			},
			expectedLists: []model.TodoList{
				{Id: 1, Title: "test1", Description: "testing1"},
				{Id: 2, Title: "test2", Description: "testing2"},
				{Id: 3, Title: "test3", Description: "testing3"},
			},
			wantErr: false,
		},
		{
			name: "Empty field",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl INNER JOIN %s ul ON (.+) WHERE (.+)", todoListsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetAll(tc.input.userId)
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

func TestTodoListSQL_GetById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoListSQL(db)

	type args struct {
		userId int
		listId int
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
				userId: 1,
				listId: 3,
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "title", "description"}).
					AddRow(3, "test", "testing")

				query := fmt.Sprintf("SELECT (.+) FROM %s tl INNER JOIN %s ul ON (.+) WHERE (.+)", todoListsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(input.userId, input.listId).WillReturnRows(rows)
			},
			expectedList: model.TodoList{Id: 3, Title: "test", Description: "testing"},
			wantErr:      false,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl INNER JOIN %s ul ON (.+) WHERE (.+)", todoListsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(0, 0)
			},
			wantErr: true,
		},
		{
			name:  "EmptyField_UserId",
			input: args{listId: 3},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl INNER JOIN %s ul ON (.+) WHERE (.+)", todoListsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(0, 3)
			},
			wantErr: true,
		},
		{
			name:  "EmptyField_ListId",
			input: args{userId: 1},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s tl INNER JOIN %s ul ON (.+) WHERE (.+)", todoListsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(1, 0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetById(tc.input.userId, tc.input.listId)
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

func TestTodoListSQL_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoListSQL(db)

	type args struct {
		listId int
		update request.UpdateTodoList
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
				listId: 1,
				update: request.UpdateTodoList{
					Title:          test.StringPointer("test"),
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("20210625:"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, *input.update.CompletionDate, input.listId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				listId: 1,
				update: request.UpdateTodoList{
					Title:          test.StringPointer("test"),
					CompletionDate: test.StringPointer("20210625:"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.CompletionDate, input.listId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitle",
			input: args{
				listId: 1,
				update: request.UpdateTodoList{
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("20210625:"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Description, *input.update.CompletionDate, input.listId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutCompletionDate",
			input: args{
				listId: 1,
				update: request.UpdateTodoList{
					Title:       test.StringPointer("test"),
					Description: test.StringPointer("testing"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET (.+) WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, input.listId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_NoUpdateFields",
			input: args{
				listId: 1,
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s tl SET WHERE (.+)", todoListsTable)
				mock.ExpectExec(query).WithArgs(input.listId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Update(tc.input.listId, tc.input.update)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoListSQL_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoListSQL(db)

	type args struct {
		listId int
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
			input: args{listId: 1},
			mockBehavior: func(input args) {
				mock.ExpectBegin()

				query1 := fmt.Sprintf("DELETE FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query1).WithArgs(input.listId).WillReturnResult(sqlmock.NewResult(1, 1))

				query2 := fmt.Sprintf("DELETE FROM %s tl WHERE (.+)", todoListsTable)
				mock.ExpectExec(query2).WithArgs(input.listId).WillReturnResult(sqlmock.NewResult(1, 1))

				query3 := fmt.Sprintf("DELETE FROM %s ul WHERE (.+)", usersListsTable)
				mock.ExpectExec(query3).WithArgs(input.listId).WillReturnResult(sqlmock.NewResult(1, 1))

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

			err := repos.Delete(tc.input.listId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
