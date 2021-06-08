package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/service"
	mockService "github.com/Lapp-coder/todo-app/internal/app/todo-app/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_createList(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, userId int, list model.TodoList)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUserId          int
		inputList            model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputBody:   `{"title":"Test","description":"testing, testing, testing..."}`,
			inputUserId: 1,
			inputList:   model.TodoList{Title: "Test", Description: "testing, testing, testing..."},
			mockBehavior: func(s *mockService.MockTodoList, userId int, list model.TodoList) {
				s.EXPECT().Create(userId, list).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"list id":1}`,
		},
		{
			name:                 "empty fields",
			inputBody:            `{}`,
			inputUserId:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userId int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "short title",
			inputBody:            `{"title": "t"}`,
			inputUserId:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userId int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "long title",
			inputBody:            `{"title": "veryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitle"}`,
			inputUserId:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userId int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "short description",
			inputBody:            `{"title": "test", "description": "t"}`,
			inputUserId:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userId int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "long description",
			inputBody:            `{"title": "test", "description": "VeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLong"}`,
			inputUserId:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userId int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:        "service failure",
			inputBody:   `{"title": "test", "description": "testing"}`,
			inputList:   model.TodoList{Title: "test", Description: "testing"},
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockTodoList, userId int, list model.TodoList) {
				s.EXPECT().Create(userId, list).Return(0, errors.New("error occurred when creating a list"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"error occurred when creating a list"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.inputUserId, tc.inputList)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.POST(
				"/api/lists",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
				},
				handler.createList)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists", bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestHandler_getAllLists(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, userId interface{}, lists []model.TodoList)

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		lists                []model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			lists: []model.TodoList{
				{Id: 1, Title: "test", Description: "testing"},
				{Id: 2, Title: "test2", Description: "testing2"},
			},
			mockBehavior: func(s *mockService.MockTodoList, userId interface{}, lists []model.TodoList) {
				s.EXPECT().GetAll(userId).Return(lists, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"lists":[{"Id":1,"Title":"test","Description":"testing"},{"Id":2,"Title":"test2","Description":"testing2"}]}`,
		},
		{
			name:        "null lists",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockTodoList, userId interface{}, lists []model.TodoList) {
				s.EXPECT().GetAll(userId).Return(lists, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"lists":null}`,
		},
		{
			name:                 "invalid user id",
			inputUserId:          "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, userId interface{}, lists []model.TodoList) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "service failure",
			inputUserId: 1,
			mockBehavior: func(s *mockService.MockTodoList, userId interface{}, lists []model.TodoList) {
				s.EXPECT().GetAll(userId).Return(nil, errors.New("failed to get all the lists"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get all the lists"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.inputUserId, tc.lists)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/lists",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
				},
				handler.getAllLists)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/lists", nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestHandler_getListById(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, list model.TodoList, userId, listId interface{})

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		list                 model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			inputParam:  1,
			list:        model.TodoList{Id: 1, Title: "test", Description: "testing"},
			mockBehavior: func(s *mockService.MockTodoList, list model.TodoList, userId, listId interface{}) {
				s.EXPECT().GetById(userId, listId).Return(list, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":{"Id":1,"Title":"test","Description":"testing"}}`,
		},
		{
			name:                 "invalid param",
			inputUserId:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, list model.TodoList, userId, listId interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:                 "invalid user id",
			inputUserId:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, list model.TodoList, userId, listId interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "service failure",
			inputUserId: 1,
			inputParam:  1,
			list:        model.TodoList{Id: 1, Title: "test", Description: "testing"},
			mockBehavior: func(s *mockService.MockTodoList, list model.TodoList, userId, listId interface{}) {
				s.EXPECT().GetById(userId, listId).Return(model.TodoList{}, errors.New("failed to get list by id"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get list by id"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.list, tc.inputUserId, tc.inputParam)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/lists/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
				},
				handler.getListById)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%v", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestHandler_updateList(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList)

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		inputBody            string
		updateList           model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing"}`,
			updateList:  model.TodoList{Title: "test", Description: "testing"},
			mockBehavior: func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {
				s.EXPECT().Update(userId, listId, list).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list update was successful"}`,
		},
		{
			name:                 "invalid param",
			inputUserId:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:                 "empty fields",
			inputUserId:          1,
			inputParam:           1,
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"update request has not values"}`,
		},
		{
			name:                 "short title",
			inputBody:            `{"title": "t"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "long title",
			inputBody:            `{"title": "veryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitle"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "short description",
			inputBody:            `{"title": "test", "description": "t"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "long description",
			inputBody:            `{"title": "test", "description": "VeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLong"}`,
			inputUserId:          1,
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "invalid user id",
			inputUserId:          "invalid",
			inputParam:           1,
			inputBody:            `{"title": "test", "description": "testing"}`,
			mockBehavior:         func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "service failure",
			inputUserId: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing"}`,
			updateList:  model.TodoList{Title: "test", Description: "testing"},
			mockBehavior: func(s *mockService.MockTodoList, userId, listId interface{}, list model.TodoList) {
				s.EXPECT().Update(userId, listId, list).Return(errors.New("failed to update the list"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to update the list"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.inputUserId, tc.inputParam, tc.updateList)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.PUT(
				"/api/lists/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
				},
				handler.updateList)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT",
				fmt.Sprintf("/api/lists/%v", tc.inputParam),
				bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}

func TestHandler_deleteList(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, dbUsersLists map[int]int, userId, listId interface{})

	testTable := []struct {
		name                 string
		inputUserId          interface{}
		inputParam           interface{}
		dbUsersLists         map[int]int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "OK",
			inputUserId:  1,
			inputParam:   1,
			dbUsersLists: map[int]int{1: 1},
			mockBehavior: func(s *mockService.MockTodoList, dbUsersLists map[int]int, userId, listId interface{}) {
				if dbUsersLists[userId.(int)] == listId {
					s.EXPECT().Delete(userId, listId).Return(nil)
					return
				}

				s.EXPECT().Delete(userId, listId).Return(errors.New("failed to delete list"))
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list deletion was successful"}`,
		},
		{
			name:                 "invalid param",
			inputUserId:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, dbUsersLists map[int]int, userId, listId interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid the param"}`,
		},
		{
			name:                 "invalid user id",
			inputUserId:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, dbUsersLists map[int]int, userId, listId interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to get user id"}`,
		},
		{
			name:        "isn't in db",
			inputUserId: 1,
			inputParam:  1,
			mockBehavior: func(s *mockService.MockTodoList, dbUsersLists map[int]int, userId, listId interface{}) {
				if dbUsersLists[userId.(int)] == listId {
					s.EXPECT().Delete(userId, listId).Return(nil)
					return
				}

				s.EXPECT().Delete(userId, listId).Return(errors.New("failed to delete list"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to delete list"}`,
		},
		{
			name:         "service failure",
			inputUserId:  1,
			inputParam:   1,
			dbUsersLists: map[int]int{1: 1},
			mockBehavior: func(s *mockService.MockTodoList, dbUsersLists map[int]int, userId, listId interface{}) {
				if dbUsersLists[userId.(int)] == listId {
					s.EXPECT().Delete(userId, listId).Return(errors.New("failed to delete list"))
					return
				}

				s.EXPECT().Delete(userId, listId).Return(errors.New("failed to delete list"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to delete list"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.dbUsersLists, tc.inputUserId, tc.inputParam)

			services := &service.Service{TodoList: todoList}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.DELETE(
				"/api/lists/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserId)
				},
				handler.deleteList)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%v", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, tc.expectedStatusCode)
			assert.Equal(t, w.Body.String(), tc.expectedResponseBody)
		})
	}
}
