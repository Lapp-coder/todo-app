package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Lapp-coder/todo-app/internal/app/service"
	mockService "github.com/Lapp-coder/todo-app/internal/app/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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

func TestHandler_getUserId(t *testing.T) {
	// Arrange
	testTable := []struct {
		name                 string
		userId               interface{}
		expectedStatusCode   int
		expectedResponseBody string
		inContext            bool
	}{
		{
			name:                 "OK",
			userId:               1,
			expectedStatusCode:   200,
			expectedResponseBody: `{"result":1}`,
			inContext:            true,
		},
		{
			name:                 "Invalid user id",
			userId:               "invalid",
			expectedStatusCode:   500,
			expectedResponseBody: `{"result":0}`,
			inContext:            true,
		},
		{
			name:                 "No in context",
			expectedStatusCode:   500,
			expectedResponseBody: `{"result":0}`,
			inContext:            false,
		},
	}

	// Act
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			// Init dependency
			services := &service.Service{}
			handler := NewHandler(services)

			// Test server
			gin.SetMode("test")
			r := gin.New()
			r.GET("/userId", func(c *gin.Context) {
				if tc.inContext {
					c.Set(userCtx, tc.userId)
				}

				id := handler.getUserId(c)
				if id == 0 {
					c.JSON(http.StatusInternalServerError, gin.H{
						"result": id,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"result": id,
				})
			})

			// Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/userId", nil)

			// Perform request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedResponseBody, w.Body.String())
		})
	}
}
