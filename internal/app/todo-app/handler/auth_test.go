package handler

import (
	"bytes"
	"errors"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/model"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/service"
	mockService "github.com/Lapp-coder/todo-app/internal/app/todo-app/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestHandler_signUp(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockAuthorization, user model.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name": "Test", "email":"test@mail.ru", "password":"testing"}`,
			inputUser: model.User{Name: "Test", Email: "test@mail.ru", Password: "testing"},
			mockBehavior: func(s *mockService.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:                 "short name",
			inputBody:            `{"name": "t", "email":"test@mail.ru", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "very long name",
			inputBody:            `{"name": "nameAndNameAndNameNameAndNameAndNameNameAndNameAndNameNameAnd", "email":"test@mail.ru", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "invalid email",
			inputBody:            `{"name": "Test", "email":"test", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "short password",
			inputBody:            `{"name": "Test", "email":"test@mail.ru", "password":"test"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "very long password",
			inputBody:            `{"name": "Test", "email":"test@mail.ru", "password":"testAndTestAndTestAndTestAndTestAndTestAndTestAndTestAndTest"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "empty fields",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:      "service failure",
			inputBody: `{"name": "Test", "email":"test@mail.ru", "password":"testing"}`,
			inputUser: model.User{Name: "Test", Email: "test@mail.ru", Password: "testing"},
			mockBehavior: func(s *mockService.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("failed to create account"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"failed to create account"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up", bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	// Arrange
	type mockBehaviour func(s *mockService.MockAuthorization, email, password string)

	testTable := []struct {
		name                 string
		inputBody            string
		inputEmail           string
		inputPassword        string
		mockBehavior         mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:          "OK",
			inputBody:     `{"email": "test@mail.ru", "password": "testing"}`,
			inputEmail:    "test@mail.ru",
			inputPassword: "testing",
			mockBehavior: func(s *mockService.MockAuthorization, email, password string) {
				s.EXPECT().GenerateToken(email, password).Return("generatedToken", nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"token":"generatedToken"}`,
		},
		{
			name:                 "invalid email",
			inputBody:            `{"email":"test", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "short password",
			inputBody:            `{"email":"test@mail.ru", "password":"test"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "very long password",
			inputBody:            `{"name": "Test", "email":"test@mail.ru", "password":"testAndTestAndTestAndTestAndTestAndTestAndTestAndTestAndTest"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:                 "empty fields",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid input body"}`,
		},
		{
			name:          "service failure",
			inputBody:     `{"email": "test@mail.ru", "password": "testing"}`,
			inputEmail:    "test@mail.ru",
			inputPassword: "testing",
			mockBehavior: func(s *mockService.MockAuthorization, email, password string) {
				s.EXPECT().GenerateToken(email, password).Return("", errors.New("incorrect email or password"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"incorrect email or password"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputEmail, tc.inputPassword)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.POST("/auth/sign-in", handler.signIn)

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in", bytes.NewBufferString(tc.inputBody))

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
