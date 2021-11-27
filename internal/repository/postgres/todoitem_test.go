package postgres

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/test"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoItemRepository(db)

	type args struct {
		listID int
		item   model.TodoItem
	}

	type mockBehavior func(input args)

	testCases := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		expectedID   int
		wantErr      bool
	}{
		{
			name: "OK_AllFields",
			input: args{
				listID: 1,
				item:   model.TodoItem{Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)
				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listID, input.item.Title, input.item.Description, input.item.CompletionDate).WillReturnRows(rows)
			},
			expectedID: 3,
			wantErr:    false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				listID: 1,
				item:   model.TodoItem{Title: "test", CompletionDate: "2021-11-21 00:00:00"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)
				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listID, input.item.Title, input.item.CompletionDate).WillReturnRows(rows)
			},
			expectedID: 3,
			wantErr:    false,
		},
		{
			name: "OK_WithoutCompletionDate",
			input: args{
				listID: 1,
				item:   model.TodoItem{Title: "test", Description: "testing"},
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)

				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listID, input.item.Title, input.item.Description).WillReturnRows(rows)
			},
			expectedID: 3,
			wantErr:    false,
		},
		{
			name: "OK_EmptyFields",
			input: args{
				listID: 1,
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id"}).AddRow(3)

				query := fmt.Sprintf("INSERT INTO %s", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listID).WillReturnRows(rows)
			},
			expectedID: 3,
			wantErr:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.Create(tc.input.listID, tc.input.item)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_GetAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoItemRepository(db)

	type args struct {
		listID int
	}

	type mockBehavior func(input args)

	testCases := []struct {
		name          string
		input         args
		mockBehavior  mockBehavior
		expectedItems []model.TodoItem
		wantErr       bool
	}{
		{
			name:  "OK",
			input: args{listID: 1},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "list_id", "title", "description", "completion_date", "done"}).
					AddRow(1, 1, "test", "testing", "2021-11-21 00:00:00", true).
					AddRow(2, 1, "test2", "testing2", "2021-11-21 00:00:00", false).
					AddRow(3, 1, "test3", "testing3", "2021-11-21 00:00:00", false)

				query := fmt.Sprintf("SELECT (.+) FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectQuery(query).WithArgs(input.listID).WillReturnRows(rows)
			},
			expectedItems: []model.TodoItem{
				{ID: 1, ListID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00", Done: true},
				{ID: 2, ListID: 1, Title: "test2", Description: "testing2", CompletionDate: "2021-11-21 00:00:00", Done: false},
				{ID: 3, ListID: 1, Title: "test3", Description: "testing3", CompletionDate: "2021-11-21 00:00:00", Done: false},
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetAll(tc.input.listID)
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

func TestTodoItemPostgres_GetByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoItemRepository(db)

	type args struct {
		userID int
		listID int
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
				userID: 1,
				listID: 3,
			},
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "list_id", "title", "description", "completion_date", "done"}).
					AddRow(1, 3, "test", "testing", "2021-11-21 00:00:00", true)

				query := fmt.Sprintf("SELECT (.+) FROM %s ti INNER JOIN %s tl ON (.+) WHERE (.+)", todoItemsTable, todoListsTable)
				mock.ExpectQuery(query).WithArgs(input.userID, input.listID).WillReturnRows(rows)
			},
			expectedItem: model.TodoItem{
				ID:             1,
				ListID:         3,
				Title:          "test",
				Description:    "testing",
				CompletionDate: "2021-11-21 00:00:00",
				Done:           true,
			},
			wantErr: false,
		},
		{
			name: "Empty fields",
			mockBehavior: func(input args) {
				rows := mock.NewRows([]string{"id", "list_id", "title", "completion_date", "description", "done"})

				query := fmt.Sprintf("SELECT (.+) FROM %s ti INNER JOIN %s tl ON (.+) WHERE (.+)", todoItemsTable, todoListsTable)
				mock.ExpectQuery(query).WithArgs(0, 0).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tc := range tesTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			got, err := repos.GetByID(tc.input.userID, tc.input.listID)
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

func TestTodoItemPostgres_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoItemRepository(db)

	type args struct {
		itemID int
		update model.UpdateTodoItem
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
				itemID: 1,
				update: model.UpdateTodoItem{
					Title:          test.StringPointer("test"),
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
					Done:           test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, *input.update.CompletionDate, *input.update.Done, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDone",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					Title:          test.StringPointer("test"),
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, *input.update.CompletionDate, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDescription",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					Title:          test.StringPointer("test"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
					Done:           test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.CompletionDate, *input.update.Done, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitle",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
					Done:           test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Description, *input.update.CompletionDate, *input.update.Done, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutDoneAndDescription",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					Title:          test.StringPointer("test"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.CompletionDate, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitleAndDone",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					Description:    test.StringPointer("testing"),
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Description, *input.update.CompletionDate, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutTitleAndDescription",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
					Done:           test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.CompletionDate, *input.update.Done, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_WithoutCompletionDate",
			input: args{
				itemID: 1,
				update: model.UpdateTodoItem{
					Title:       test.StringPointer("test"),
					Description: test.StringPointer("testing"),
					Done:        test.BoolPointer(true),
				},
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET (.+) WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(
					*input.update.Title, *input.update.Description, *input.update.Done, input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "OK_NoUpdateFields",
			input: args{
				itemID: 1,
			},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("UPDATE %s ti SET WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
	}

	for _, tc := range tesTable {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Update(tc.input.itemID, tc.input.update)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestTodoItemPostgres_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error occurred when opening a connection to the stub database: %s", err.Error())
	}

	db := sqlx.NewDb(mockDB, "sqlmock")
	repos := NewTodoItemRepository(db)

	type args struct {
		itemID int
	}

	type mockBehavior func(input args)

	testCases := []struct {
		name         string
		input        args
		mockBehavior mockBehavior
		wantErr      bool
	}{
		{
			name:  "OK",
			input: args{itemID: 1},
			mockBehavior: func(input args) {
				query := fmt.Sprintf("DELETE FROM %s ti WHERE (.+)", todoItemsTable)
				mock.ExpectExec(query).WithArgs(input.itemID).WillReturnResult(sqlmock.NewResult(1, 1))
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(tc.input)

			err := repos.Delete(tc.input.itemID)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
