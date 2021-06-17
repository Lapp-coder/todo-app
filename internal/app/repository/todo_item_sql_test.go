package repository

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/Lapp-coder/todo-app/internal/app/request"
	"github.com/Lapp-coder/todo-app/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoItemSQL_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoItemSQL(db)

	type args struct {
		listId int
		item   model.TodoItem
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		expectedId   int
		wantErr      bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				listId: 1,
				item:   model.TodoItem{Title: "test", Description: "testing"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)

				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listId, input.item.Title, input.item.Description).WillReturnRows(rows)
			},
			expectedId: 3,
			wantErr:    false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				listId: 1,
				item:   model.TodoItem{Title: "test"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)

				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listId, input.item.Title, "").WillReturnRows(rows)
			},
			expectedId: 3,
			wantErr:    false,
		},
		{
			name: "OK_Emptyfields",
			input: args{
				listId: 1,
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)

				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listId, "", "").WillReturnRows(rows)
			},
			expectedId: 3,
			wantErr:    false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.Create(tc.input.listId, tc.input.item)
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

func TestTodoItemSQL_GetAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoItemSQL(db)

	type args struct {
		listId int
	}

	type mockBehavior func(input args)

	testTable := []struct {
		name          string
		input         args
		mockBehavior  mockBehavior
		expectedItems []model.TodoItem
		wantErr       bool
	}{
		{
			name:  "OK",
			input: args{listId: 1},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "list_id", "title", "description", "done"}).
					AddRow(1, 1, "test", "testing", true).
					AddRow(2, 1, "test2", "testing2", false).
					AddRow(3, 1, "test3", "testing3", false)

				query := fmt.Sprintf("SELECT (.+) FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listId).WillReturnRows(rows)
			},
			expectedItems: []model.TodoItem{
				{1, 1, "test", "testing", true},
				{2, 1, "test2", "testing2", false},
				{3, 1, "test3", "testing3", false},
			},
			wantErr: false,
		},
		{
			name: "Empty field",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("SELECT (.+) FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(0)
			},
			wantErr: true,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetAll(tc.input.listId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedItems, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemSQL_GetById(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoItemSQL(db)

	type args struct {
		userId int
		listId int
	}

	type mockBehavior func(input args)

	tesTable := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		expectedItem model.TodoItem
		wantErr      bool
	}{
		{
			name: "OK",
			input: args{
				userId: 1,
				listId: 3,
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "list_id", "title", "description", "done"}).
					AddRow(1, 3, "test", "testing", true)

				query := fmt.Sprintf("SELECT (.+) FROM %s ti INNER JOIN %s ul ON (.+) WHERE (.+)", todoItemsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(input.userId, input.listId).WillReturnRows(rows)
			},
			expectedItem: model.TodoItem{
				Id:          1,
				ListId:      3,
				Title:       "test",
				Description: "testing",
				Done:        true,
			},
			wantErr: false,
		},
		{
			name: "Empty fields",
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "list_id", "title", "description", "done"})

				query := fmt.Sprintf("SELECT (.+) FROM %s ti INNER JOIN %s ul ON (.+) WHERE (.+)", todoItemsTable, usersListsTable)
				mock.ExpectQuery(query).WithArgs(0, 0).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tc := range tesTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetById(tc.input.userId, tc.input.listId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedItem, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemSQL_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoItemSQL(db)

	type args struct {
		itemId int
		update request.UpdateTodoItem
	}

	type mockBehavior func(input args)

	tesTable := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Title:       test.StringPointer("test"),
					Description: test.StringPointer("testing"),
					Done:        test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, *input.update.Done, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDone",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Title:       test.StringPointer("test"),
					Description: test.StringPointer("testing"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Title: test.StringPointer("test"),
					Done:  test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Done, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitle",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Description: test.StringPointer("testing"),
					Done:        test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Description, *input.update.Done, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDoneAndDescription",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Title: test.StringPointer("test"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitleAndDone",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Description: test.StringPointer("testing"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Description, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitleAndDescription",
			input: args{
				itemId: 1,
				update: request.UpdateTodoItem{
					Done: test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Done, input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_NoUpdateFields",
			input: args{
				itemId: 1,
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tc := range tesTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Update(tc.input.itemId, tc.input.update)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemSQL_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")

	repos := NewTodoItemSQL(db)

	type args struct {
		itemId int
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
			input: args{itemId: 1},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("DELETE FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(input.itemId).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_EmptyField",
			mockBehavior: func(input args) {
				query := fmt.Sprintf("DELETE FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(0).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Delete(tc.input.itemId)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
