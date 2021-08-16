package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockTokenRepository struct{}

func (tr *mockTokenRepository) IsTokenValid(ctx context.Context, token string) (bool, error) {
	if token == "321" {
		return false, nil
	}

	if token == "error" {
		return false, fmt.Errorf("Error deliberately caused due to wrong token")
	}

	return true, nil
}

func (tr *mockTokenRepository) InsertToken(ctx context.Context, token, created string) error {
	return nil
}

func TestAuthMiddleware_Authenticate(t *testing.T) {
	t.Run("All OK", func(t *testing.T) {
		// Given
		request, requestRecorder, engine := newRouter("GET", "/test", nil)
		request.Header.Set("Authorization", "Bearer 123")

		// When
		engine.Use(NewAuthMiddleware(&mockTokenRepository{}).Authenticate)
		engine.GET("/test", func(c *gin.Context) {
			tokenIntf, _ := c.Get("token")
			if tokenIntf.(string) != "123" {
				c.String(http.StatusOK, "FAIL")
			}

			c.String(http.StatusOK, "PASS")
		})
		engine.ServeHTTP(requestRecorder, request)

		body, err := ioutil.ReadAll(requestRecorder.Body)
		if err != nil {
			t.Fatalf("Failed to ready body content: %v", err)
		}

		// Then
		assert.Equal(t, http.StatusOK, requestRecorder.Code)
		assert.Equal(t, "PASS", string(body))
	})

	t.Run("No auth", func(t *testing.T) {
		// Given
		request, requestRecorder, engine := newRouter("GET", "/test", nil)

		// When
		engine.Use(NewAuthMiddleware(&mockTokenRepository{}).Authenticate)
		engine.GET("/test")
		engine.ServeHTTP(requestRecorder, request)

		// Then
		assert.Equal(t, http.StatusUnauthorized, requestRecorder.Code)
	})

	t.Run("Invalid token", func(t *testing.T) {
		// Given
		request, requestRecorder, engine := newRouter("GET", "/test", nil)
		request.Header.Set("Authorization", "Bearer 321")

		// When
		engine.Use(NewAuthMiddleware(&mockTokenRepository{}).Authenticate)
		engine.GET("/test")
		engine.ServeHTTP(requestRecorder, request)

		// Then
		assert.Equal(t, http.StatusUnauthorized, requestRecorder.Code)
	})

	t.Run("Repo error", func(t *testing.T) {
		// Given
		request, requestRecorder, engine := newRouter("GET", "/test", nil)
		request.Header.Set("Authorization", "Bearer error")

		// When
		engine.Use(NewAuthMiddleware(&mockTokenRepository{}).Authenticate)
		engine.GET("/test")
		engine.ServeHTTP(requestRecorder, request)

		// Then
		assert.Equal(t, http.StatusInternalServerError, requestRecorder.Code)
	})
}

func newRouter(method, path string, body io.Reader) (*http.Request, *httptest.ResponseRecorder, *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	req := httptest.NewRequest(method, path, body)
	rr := httptest.NewRecorder()
	_, eng := gin.CreateTestContext(rr)

	return req, rr, eng
}
