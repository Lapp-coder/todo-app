package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Lapp-coder/todo-app/internal/app/model"
	"github.com/Lapp-coder/todo-app/internal/app/request"
	"github.com/Lapp-coder/todo-app/internal/app/service"
	mockService "github.com/Lapp-coder/todo-app/internal/app/service/mocks"
	"github.com/Lapp-coder/todo-app/test"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createItem(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem)

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		inputBody            string
		item                 model.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK_AllFields",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "20210626:", "done": false}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: "20210626 ", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {
				s.EXPECT().Create(userId, listId, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item id":1}`,
		},
		{
			name:        "OK_WithoutDone",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description":"testing", "completion_date": "20210626:"}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: "20210626 "},
			mockBehavior: func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {
				s.EXPECT().Create(userId, listId, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item id":1}`,
		},
		{
			name:        "OK_WithoutDescription",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "completion_date":"20210626:"}`,
			item:        model.TodoItem{Title: "test", CompletionDate: "20210626 "},
			mockBehavior: func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {
				s.EXPECT().Create(userId, listId, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item id":1}`,
		},
		{
			name:        "OK_WithoutCompletionDate",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description":"testing"}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: test.GetTimeNow()},
			mockBehavior: func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {
				s.EXPECT().Create(userId, listId, item).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"item id":1}`,
		},
		{
			name:                 "Invalid param",
			inputUserId:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:                 "Empty fields",
			inputUserId:          1,
			inputParam:           1,
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Short title",
			inputBody:            `{"title": "t"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Long title",
			inputBody:            `{"title": "veryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitle"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Short description",
			inputBody:            `{"title": "test", "description": "t"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Long description",
			inputBody:            `{"title": "test", "description": "VeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLong"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:        "Invalid user id",
			inputUserId: "invalid",
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "Service failure",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description":"testing", "completion_date": "20210626:"}`,
			item:        model.TodoItem{Title: "test", Description: "testing", CompletionDate: "20210626 "},
			mockBehavior: func(s *mockService.MockTodoItem, userId, listId interface{}, item model.TodoItem) {
				s.EXPECT().Create(userId, listId, item).Return(0, errors.New("failed to create item"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to create item"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserId, tc.inputParam, tc.item)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.POST(
				"/api/lists/:id/items",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
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
	type mockBehavior func(s *mockService.MockTodoItem, items []model.TodoItem, userId, listId interface{})

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		items                []model.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			inputParam:  1,
			items: []model.TodoItem{
				{Id: 1, ListId: 1, Title: "test", Description: "testing", CompletionDate: "20210626", Done: false},
				{Id: 2, ListId: 1, Title: "test2", Description: "testing2", CompletionDate: "20210626", Done: true},
			},
			mockBehavior: func(s *mockService.MockTodoItem, items []model.TodoItem, userId, listId interface{}) {
				s.EXPECT().GetAll(userId, listId).Return(items, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"items":[{"Id":1,"ListId":1,"Title":"test","Description":"testing","CompletionDate":"20210626","Done":false},{"Id":2,"ListId":1,"Title":"test2","Description":"testing2","CompletionDate":"20210626","Done":true}]}`,
		},
		{
			name:        "Invalid param",
			inputUserId: 1,
			inputParam:  "invalid",
			mockBehavior: func(s *mockService.MockTodoItem, items []model.TodoItem, userId, listId interface{}) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:        "Invalid user id",
			inputUserId: "invalid",
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, items []model.TodoItem, userId, listId interface{}) {
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "Service failure",
			inputUserId: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, items []model.TodoItem, userId, listId interface{}) {
				s.EXPECT().GetAll(userId, listId).Return(nil, errors.New("failed to get all items"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get all items"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.items, tc.inputUserId, tc.inputParam)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/lists/:id/items",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
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

func TestHandler_getItemById(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoItem, userId, itemId interface{}, item model.TodoItem)

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		item                 model.TodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			inputParam:  1,
			item:        model.TodoItem{Id: 1, ListId: 1, Title: "test", Description: "testing", CompletionDate: "20210626 ", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, item model.TodoItem) {
				s.EXPECT().GetById(userId, itemId).Return(item, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"item":{"Id":1,"ListId":1,"Title":"test","Description":"testing","CompletionDate":"20210626 ","Done":false}}`,
		},
		{
			name:        "Invalid param",
			inputUserId: 1,
			inputParam:  "invalid",
			item:        model.TodoItem{Id: 1, ListId: 1, Title: "test", Description: "testing", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, item model.TodoItem) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:        "Invalid user id",
			inputUserId: "invalid",
			inputParam:  1,
			item:        model.TodoItem{Id: 1, ListId: 1, Title: "test", Description: "testing", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, item model.TodoItem) {

			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "Service failure",
			inputUserId: 1,
			inputParam:  1,
			item:        model.TodoItem{Id: 1, ListId: 1, Title: "test", Description: "testing", CompletionDate: "20210626", Done: false},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, item model.TodoItem) {
				s.EXPECT().GetById(userId, itemId).Return(item, errors.New("failed to get item by id"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get item by id"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserId, tc.inputParam, tc.item)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/items/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
				},
				handler.getItemById)

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
	type mockBehavior func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem)

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		inputBody            string
		update               request.UpdateTodoItem
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK_AllFields",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "20210626:", "done": true}`,
			update: request.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("20210626 "),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {
				s.EXPECT().Update(userId, itemId, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutDone",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "20210626:"}`,
			update: request.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("20210626 "),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {
				s.EXPECT().Update(userId, itemId, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutDescriptionAndDone",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "completion_date": "20210626:"}`,
			update: request.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				CompletionDate: test.StringPointer("20210626 "),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {
				s.EXPECT().Update(userId, itemId, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutTitleAndDescription",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"completion_date": "20210626:", "done": true}`,
			update: request.UpdateTodoItem{
				CompletionDate: test.StringPointer("20210626 "),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {
				s.EXPECT().Update(userId, itemId, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:        "OK_WithoutCompletionDate",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "done": true}`,
			update: request.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer(test.GetTimeNow()),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {
				s.EXPECT().Update(userId, itemId, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item update was successful"}`,
		},
		{
			name:                 "Invalid param",
			inputUserId:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:                 "Empty fields",
			inputUserId:          1,
			inputParam:           1,
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"update request has not values"}`,
		},
		{
			name:                 "Short title",
			inputBody:            `{"title": "t"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Long title",
			inputBody:            `{"title": "veryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitle"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Short description",
			inputBody:            `{"title": "test", "description": "t"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Long description",
			inputBody:            `{"title": "test", "description": "VeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLong"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "Invalid user id",
			inputUserId:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "Service failure",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "20210626:", "done": true}`,
			update: request.UpdateTodoItem{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("20210626 "),
				Done:           test.BoolPointer(true),
			},
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}, update request.UpdateTodoItem) {
				s.EXPECT().Update(userId, itemId, update).Return(errors.New("failed to update item"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to update item"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserId, tc.inputParam, tc.update)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.PUT(
				"/api/items/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
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
	type mockBehavior func(s *mockService.MockTodoItem, userId, itemId interface{})

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}) {
				s.EXPECT().Delete(userId, itemId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the item deletion was successful"}`,
		},
		{
			name:                 "Invalid param",
			inputUserId:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:                 "Invalid user id",
			inputUserId:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoItem, userId, itemId interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "Service failure",
			inputUserId: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoItem, userId, itemId interface{}) {
				s.EXPECT().Delete(userId, itemId).Return(errors.New("failed to delete item"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to delete item"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoItem := mockService.NewMockTodoItem(c)
			tc.mockBehavior(todoItem, tc.inputUserId, tc.inputParam)

			services := &service.Service{TodoItem: todoItem}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.DELETE(
				"/api/items/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
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
