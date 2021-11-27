package handler

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Lapp-coder/todo-app/test"

	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/service"
	mockService "github.com/Lapp-coder/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_createList(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, userID int, list model.TodoList)

	testCases := []struct {
		name                 string
		inputBody            string
		inputUserID          int
		inputList            model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK_AllFields",
			inputBody:   `{"title":"Test","description":"testing, testing, testing...","completion_date":"2021-11-21 00:00:00"}`,
			inputUserID: 1,
			inputList:   model.TodoList{Title: "Test", Description: "testing, testing, testing...", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoList, userID int, list model.TodoList) {
				s.EXPECT().Create(userID, list).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"list id":1}`,
		},
		{
			name:        "OK_WithoutDescription",
			inputBody:   `{"title":"Test","completion_date":"2021-11-21 00:00:00"}`,
			inputUserID: 1,
			inputList:   model.TodoList{Title: "Test", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoList, userID int, list model.TodoList) {
				s.EXPECT().Create(userID, list).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"list id":1}`,
		},
		{
			name:        "OK_WithoutCompletionDate",
			inputBody:   `{"title":"Test","description":"testing"}`,
			inputUserID: 1,
			inputList:   model.TodoList{Title: "Test", Description: "testing"},
			mockBehavior: func(s *mockService.MockTodoList, userID int, list model.TodoList) {
				s.EXPECT().Create(userID, list).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"list id":1}`,
		},
		{
			name:                 "Empty fields",
			inputBody:            `{}`,
			inputUserID:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userID int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Short title",
			inputBody:            `{"title": "t"}`,
			inputUserID:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userID int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long title",
			inputBody:            `{"title": "veryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitleVeryLongTitle"}`,
			inputUserID:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userID int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long description",
			inputBody:            `{"title": "test", "description": "VeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLongDescriptionVeryLong"}`,
			inputUserID:          1,
			mockBehavior:         func(s *mockService.MockTodoList, userID int, list model.TodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:        "Service failure",
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00"}`,
			inputList:   model.TodoList{Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			inputUserID: 1,
			mockBehavior: func(s *mockService.MockTodoList, userID int, list model.TodoList) {
				s.EXPECT().Create(userID, list).Return(0, service.ErrFailedToCreateList)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToCreateList.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.inputUserID, tc.inputList)

			services := &service.Service{TodoList: todoList}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.POST(
				"/api/lists",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.createList)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/lists", bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getAllLists(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, userID interface{}, lists []model.TodoList)

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		lists                []model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserID: 1,
			lists: []model.TodoList{
				{ID: 1, UserID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
				{ID: 2, UserID: 1, Title: "test2", Description: "testing2", CompletionDate: "2021-11-21 00:00:00"},
			},
			mockBehavior: func(s *mockService.MockTodoList, userID interface{}, lists []model.TodoList) {
				s.EXPECT().GetAll(userID).Return(lists, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"lists":[{"id":1,"user_id":1,"title":"test","description":"testing","completion_date":"2021-11-21 00:00:00"},{"id":2,"user_id":1,"title":"test2","description":"testing2","completion_date":"2021-11-21 00:00:00"}]}`,
		},
		{
			name:        "No lists",
			inputUserID: 1,
			mockBehavior: func(s *mockService.MockTodoList, userID interface{}, lists []model.TodoList) {
				s.EXPECT().GetAll(userID).Return(lists, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"lists":null}`,
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, userID interface{}, lists []model.TodoList) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			mockBehavior: func(s *mockService.MockTodoList, userID interface{}, lists []model.TodoList) {
				s.EXPECT().GetAll(userID).Return(nil, service.ErrFailedToGetAllLists)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToGetAllLists.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.inputUserID, tc.lists)

			services := &service.Service{TodoList: todoList}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/lists",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.getAllLists)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/api/lists", nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_getListByID(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, list model.TodoList, userID, listID interface{})

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		list                 model.TodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			inputUserID: 1,
			inputParam:  1,
			list:        model.TodoList{ID: 1, UserID: 1, Title: "test", Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoList, list model.TodoList, userID, listID interface{}) {
				s.EXPECT().GetByID(userID, listID).Return(list, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"list":{"id":1,"user_id":1,"title":"test","description":"testing","completion_date":"2021-11-21 00:00:00"}}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, list model.TodoList, userID, listID interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, list model.TodoList, userID, listID interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			list:        model.TodoList{ID: 1, Title: "test", UserID: 1, Description: "testing", CompletionDate: "2021-11-21 00:00:00"},
			mockBehavior: func(s *mockService.MockTodoList, list model.TodoList, userID, listID interface{}) {
				s.EXPECT().GetByID(userID, listID).Return(model.TodoList{}, service.ErrFailedToGetListByID)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToGetListByID.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.list, tc.inputUserID, tc.inputParam)

			services := &service.Service{TodoList: todoList}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET(
				"/api/lists/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.getListByID)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/api/lists/%v", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_updateList(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList)

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		inputBody            string
		updateList           model.UpdateTodoList
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK_AllFields",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00"}`,
			updateList: model.UpdateTodoList{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
			},
			mockBehavior: func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {
				s.EXPECT().Update(userID, listID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list update was successful"}`,
		},
		{
			name:        "OK_WithoutDescription",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "completion_date": "2021-11-21 00:00:00"}`,
			updateList: model.UpdateTodoList{
				Title:          test.StringPointer("test"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
			},
			mockBehavior: func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {
				s.EXPECT().Update(userID, listID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list update was successful"}`,
		},
		{
			name:        "OK_WithoutTitle",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"description": "testing", "completion_date": "2021-11-21 00:00:00"}`,
			updateList: model.UpdateTodoList{
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
			},
			mockBehavior: func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {
				s.EXPECT().Update(userID, listID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list update was successful"}`,
		},
		{
			name:        "OK_WithoutCompletionDate",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing"}`,
			updateList: model.UpdateTodoList{
				Title:       test.StringPointer("test"),
				Description: test.StringPointer("testing"),
			},
			mockBehavior: func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {
				s.EXPECT().Update(userID, listID, update).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list update was successful"}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Empty fields",
			inputUserID:          1,
			inputParam:           1,
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			inputBody:            `{"title": "test", "description": "testing"}`,
			mockBehavior:         func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:        "Service failure",
			inputUserID: 1,
			inputParam:  1,
			inputBody:   `{"title": "test", "description": "testing", "completion_date": "2021-11-21 00:00:00"}`,
			updateList: model.UpdateTodoList{
				Title:          test.StringPointer("test"),
				Description:    test.StringPointer("testing"),
				CompletionDate: test.StringPointer("2021-11-21 00:00:00"),
			},
			mockBehavior: func(s *mockService.MockTodoList, userID, listID interface{}, update model.UpdateTodoList) {
				s.EXPECT().Update(userID, listID, update).Return(service.ErrFailedToUpdateList)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToUpdateList.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.inputUserID, tc.inputParam, tc.updateList)

			services := &service.Service{TodoList: todoList}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.PUT(
				"/api/lists/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
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
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteList(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockTodoList, dbUsersLists map[int]int, userID, listID interface{})

	testCases := []struct {
		name                 string
		inputUserID          interface{}
		inputParam           interface{}
		dbUsersLists         map[int]int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:         "OK",
			inputUserID:  1,
			inputParam:   1,
			dbUsersLists: map[int]int{1: 1},
			mockBehavior: func(s *mockService.MockTodoList, dbUsersLists map[int]int, userID, listID interface{}) {
				if dbUsersLists[userID.(int)] == listID {
					s.EXPECT().Delete(userID, listID).Return(nil)
					return
				}

				s.EXPECT().Delete(userID, listID).Return(service.ErrFailedToDeleteList)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":"the list deletion was successful"}`,
		},
		{
			name:                 "Invalid param",
			inputUserID:          1,
			inputParam:           "invalid",
			mockBehavior:         func(s *mockService.MockTodoList, dbUsersLists map[int]int, userID, listID interface{}) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidParamID.Error()),
		},
		{
			name:                 "Invalid user id",
			inputUserID:          "invalid",
			inputParam:           1,
			mockBehavior:         func(s *mockService.MockTodoList, dbUsersLists map[int]int, userID, listID interface{}) {},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errFailedToGetUserID.Error()),
		},
		{
			name:         "Service failure",
			inputUserID:  1,
			inputParam:   1,
			dbUsersLists: map[int]int{1: 1},
			mockBehavior: func(s *mockService.MockTodoList, dbUsersLists map[int]int, userID, listID interface{}) {
				if dbUsersLists[userID.(int)] == listID {
					s.EXPECT().Delete(userID, listID).Return(service.ErrFailedToDeleteList)
					return
				}

				s.EXPECT().Delete(userID, listID).Return(service.ErrFailedToDeleteList)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrFailedToDeleteList.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			todoList := mockService.NewMockTodoList(c)
			tc.mockBehavior(todoList, tc.dbUsersLists, tc.inputUserID, tc.inputParam)

			services := &service.Service{TodoList: todoList}
			handler := New(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.DELETE(
				"/api/lists/:id",
				func(c *gin.Context) {
					c.Set(userCtx, tc.inputUserID)
				},
				handler.deleteList)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/lists/%v", tc.inputParam), nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
