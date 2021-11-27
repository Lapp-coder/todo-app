package handler

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/Lapp-coder/todo-app/internal/model"
	"github.com/Lapp-coder/todo-app/internal/service"
	mockService "github.com/Lapp-coder/todo-app/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockAuthorization, user model.User)

	testCases := []struct {
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
			name:                 "Short name",
			inputBody:            `{"name": "t", "email":"test@mail.ru", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long name",
			inputBody:            `{"name": "nameAndNameAndNameNameAndNameAndNameNameAndNameAndNameNameAnd", "email":"test@mail.ru", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Invalid email",
			inputBody:            `{"name": "Test", "email":"test", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Short password",
			inputBody:            `{"name": "Test", "email":"test@mail.ru", "password":"test"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long password",
			inputBody:            `{"name": "Test", "email":"test@mail.ru", "password":"testAndTestAndTestAndTestAndTestAndTestAndTestAndTestAndTest"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Empty fields",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockAuthorization, user model.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:      "Service failure",
			inputBody: `{"name": "Test", "email":"test@mail.ru", "password":"testing"}`,
			inputUser: model.User{Name: "Test", Email: "test@mail.ru", Password: "testing"},
			mockBehavior: func(s *mockService.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(user).Return(0, service.ErrIncorrectEmailOrPassword)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrIncorrectEmailOrPassword.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputUser)

			services := &service.Service{Authorization: auth}
			handler := New(services)

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

	testCases := []struct {
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
			name:                 "Invalid email",
			inputBody:            `{"email":"test", "password":"testing"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Short password",
			inputBody:            `{"email":"test@mail.ru", "password":"test"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Long password",
			inputBody:            `{"email":"test@mail.ru", "password":"testAndTestAndTestAndTestAndTestAndTestAndTestAndTestAndTest"}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:                 "Empty fields",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockAuthorization, email, password string) {},
			expectedStatusCode:   400,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, errInvalidInputBody.Error()),
		},
		{
			name:          "Service failure",
			inputBody:     `{"email": "test@mail.ru", "password": "testing"}`,
			inputEmail:    "test@mail.ru",
			inputPassword: "testing",
			mockBehavior: func(s *mockService.MockAuthorization, email, password string) {
				s.EXPECT().GenerateToken(email, password).Return("", service.ErrIncorrectEmailOrPassword)
			},
			expectedStatusCode:   500,
			expectedResponseBody: fmt.Sprintf(`{"error":"%s"}`, service.ErrIncorrectEmailOrPassword.Error()),
		},
	}

	// Act
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.inputEmail, tc.inputPassword)

			services := &service.Service{Authorization: auth}
			handler := New(services)

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
