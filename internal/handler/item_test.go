package handler

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/service"
	mockService "github.com/Lapp-coder/todo-app/internal/service/mocks"
	"github.com/Lapp-coder/todo-app/test"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createItem(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem)

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		inputBody            string
		item                 model.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK_AllFields",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00", "done": false}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {
				s.EXPECT().Create(userID, listID, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item_id":1}`,
		},
		{
			name:        "OK_WithoutDone",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description":"testing", "completion_date": "2021-11-21 00:00:00"}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {
				s.EXPECT().Create(userID, listID, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item_id":1}`,
		},
		{
			name:        "OK_WithoutDescription",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "completion_date": "2021-11-21 00:00:00"}`,
			item:        model.TodoItem{Title: "test", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {
				s.EXPECT().Create(userID, listID, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item_id":1}`,
		},
		{
			name:        "OK_WithoutCompletionDate",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description":"testing"}`,
			item:        model.TodoItem{Title: "test", Description: "testing"},
			mockBehavior: func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {
				s.EXPECT().Create(userID, listID, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item_id":1}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Empty fields",
			inputUserID:          1,
			inputParam:           1,
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Short title",
			inputBody:            `{"title": "t"}`,
			inputUserID:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long title",
			inputBody:            `{"title": "veryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitle"}`,
			inputUserID:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long description",
			inputBody:            `{"title": "test", "description": "VeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLong"}`,
			inputUserID:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:        "Invalid user id",
			inputUserID: "invalid",
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description":"testing", "completion_date": "2021-11-21 00:00:00"}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoItem, userID, listID interface{}, item model.TodoItem) {
				s.EXPECT().Create(userID, listID, item).Return(0, service.ErrFailedToCreateItem)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToCreateItem.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserID, tc.inputParam, tc.item)

			services := &service.Service{TodoItem: todoItem}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.POST(
				"/api/lists/:id/items",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.createItem)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"POST",
				fmt.Sprintf("/api/lists/%v/items", tc.inputParam),
				bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAllItems(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoItem, items []model.TodoItem, userID, listID interface{})

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		items                []model.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserID: 1,
			inputParam:  1,
			items: []model.TodoItem{
				{ID: 1, ListID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00", Done: false},
				{ID: 2, ListID: 1, Title: "test2", Description: "testing2", CompletionDate: "2021-11-21 00:00:00", Done: true},
			},
			mockBehavior: func(s *mockService.MockTodoItem, items []model.TodoItem, userID, listID interface{}) {
				s.EXPECT().GetAll(userID, listID).Return(items, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"items":[{"id":1,"list_id":1,"title":"test","description":"testing","completion_date":"2021-11-21 00:00:00","done":false},{"id":2,"list_id":1,"title":"test2","description":"testing2","completion_date":"2021-11-21 00:00:00","done":true}]}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, items []model.TodoItem, userID, listID interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, items []model.TodoItem, userID, listID interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, items []model.TodoItem, userID, listID interface{}) {
				s.EXPECT().GetAll(userID, listID).Return(nil, service.ErrFailedToGetAllItems)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToGetAllItems.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.items, tc.inputUserID, tc.inputParam)

			services := &service.Service{TodoItem: todoItem}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/lists/:id/items",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.getAllItems)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%v/items", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getItemByID(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoItem, userID, itemID interface{}, item model.TodoItem)

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		item                 model.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserID: 1,
			inputParam:  1,
			item:        model.TodoItem{ID: 1, ListID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, item model.TodoItem) {
				s.EXPECT().GetByID(userID, itemID).Return(item, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"item":{"id":1,"list_id":1,"title":"test","description":"testing","completion_date":"2021-11-21 00:00:00","done":false}}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			item:                 model.TodoItem{ID: 1, ListID: 1, Title: "test", Description: "testing", Done: false},
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			item:                 model.TodoItem{ID: 1, ListID: 1, Title: "test", Description: "testing", Done: false},
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}, item model.TodoItem) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			item:        model.TodoItem{ID: 1, ListID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, item model.TodoItem) {
				s.EXPECT().GetByID(userID, itemID).Return(item, service.ErrFailedToGetItemByID)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToGetItemByID.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserID, tc.inputParam, tc.item)

			services := &service.Service{TodoItem: todoItem}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/items/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.getItemByID)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/items/%v", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateItem(t *testing.T) {
	// Assert
	type mockBehavior func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem)

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		inputBody            string
		update               model.UpdateTodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK_AllFields",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00", "done": true}`,
			update: model.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {
				s.EXPECT().Update(userID, itemID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutDone",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00"}`,
			update: model.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {
				s.EXPECT().Update(userID, itemID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutDescriptionAndDone",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "completion_date": "2021-11-21 00:00:00"}`,
			update: model.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {
				s.EXPECT().Update(userID, itemID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutTitleAndDescription",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"completion_date": "2021-11-21 00:00:00", "done": true}`,
			update: model.UpdateTodoItem{
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {
				s.EXPECT().Update(userID, itemID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutCompletionDate",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "done": true}`,
			update: model.UpdateTodoItem{
				Title:       test.StringPointer("test"),
				Description: test.StringPointer("testing"),
				Done:        test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {
				s.EXPECT().Update(userID, itemID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Empty fields",
			inputUserID:          1,
			inputParam:           1,
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00", "done": true}`,
			update: model.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}, update model.UpdateTodoItem) {
				s.EXPECT().Update(userID, itemID, update).Return(service.ErrFailedToUpdateItem)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToUpdateItem.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserID, tc.inputParam, tc.update)

			services := &service.Service{TodoItem: todoItem}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.PUT(
				"/api/items/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.updateItem)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/api/items/%v", tc.inputParam), bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteItem(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoItem, userID, itemID interface{})

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserID: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}) {
				s.EXPECT().Delete(userID, itemID).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item deletion was successful"}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userID, itemID interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, userID, itemID interface{}) {
				s.EXPECT().Delete(userID, itemID).Return(service.ErrFailedToDeleteItem)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToDeleteItem.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserID, tc.inputParam)

			services := &service.Service{TodoItem: todoItem}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.DELETE(
				"/api/items/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.deleteItem)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/items/%v", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
