package handler

import (
	"errors"
	"github.com/Lapp-coder/todo-app/internal/app/todo-app/service"
	mockService "github.com/Lapp-coder/todo-app/internal/app/todo-app/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_userAuthentication(t *testing.T) {
	// Arrange
	type mockBehavior func(s *mockService.MockAuthorization, token string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Empty header name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Invalid header name",
			headerName:           "invalid",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Empty header values",
			headerName:           "Authorization",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"empty auth header"}`,
		},
		{
			name:                 "Invalid bearer",
			headerName:           "Authorization",
			headerValue:          "NoBearer token",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"invalid auth header"}`,
		},
		{
			name:                 "Empty token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"token is empty"}`,
		},
		{
			name:        "Service failure",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, errors.New("failed to parse token"))
			},
			expectedStatusCode:   401,
			expectedResponseBody: `{"error":"failed to parse token"}`,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mockService.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.token)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET("/protected", handler.userAuthentication, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, strconv.Itoa(id.(int)))
			})

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set(tc.headerName, tc.headerValue)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
