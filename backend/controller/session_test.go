package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kristofaranyos/tech-challenge-time/repository"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockSessionRepository struct{}

func (sr *mockSessionRepository) StartSession(ctx context.Context, token, createdAt string) (string, error) {
	if token == "error" {
		return "", fmt.Errorf("Error deliberately caused due to wrong token")
	}

	return "321", nil
}

func (sr *mockSessionRepository) StopSession(ctx context.Context, id, stoppedAt string) error {
	return nil
}

func (sr *mockSessionRepository) GetSession(ctx context.Context, id string) (*repository.Session, error) {
	return nil, nil
}

func (sr *mockSessionRepository) ListSessions(ctx context.Context, token, from, to string) ([]*repository.Session, error) {
	return nil, nil
}

func (sr *mockSessionRepository) SetName(ctx context.Context, id, name string) error {
	return nil
}

func TestSessionController_Start(t *testing.T) {
	t.Run("All OK", func(t *testing.T) {
		// Given
		request, requestRecorder, engine := newRouter("POST", "/api/v1/start", nil, "123")

		// When
		engine.POST("/api/v1/start", NewSessionController(&mockSessionRepository{}).Start)
		engine.ServeHTTP(requestRecorder, request)

		body, err := ioutil.ReadAll(requestRecorder.Body)
		if err != nil {
			t.Fatalf("Failed to ready body content: %v", err)
		}

		// Then
		assert.Equal(t, http.StatusOK, requestRecorder.Code)
		assert.Equal(t, "{\"sessionId\":\"321\"}", string(body))
	})

	t.Run("Repo error", func(t *testing.T) {
		// Given
		request, requestRecorder, engine := newRouter("POST", "/api/v1/start", nil, "error")

		// When
		engine.POST("/api/v1/start", NewSessionController(&mockSessionRepository{}).Start)
		engine.ServeHTTP(requestRecorder, request)

		// Then
		assert.Equal(t, http.StatusInternalServerError, requestRecorder.Code)
	})
}

// Rest of the endpoints tested in the similar manner... you get it
// List endpoint needs very careful testing to ensure there aren't problems with boundary date values
// Or even better, the date calculation can be refactored out as a separate function and can be tested independently

func newRouter(method, path string, body io.Reader, token string) (*http.Request, *httptest.ResponseRecorder, *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)
	req := httptest.NewRequest(method, path, body)
	rr := httptest.NewRecorder()
	_, eng := gin.CreateTestContext(rr)

	if token != "" {
		eng.Use(func(c *gin.Context) {
			c.Set("token", token)
		})
	}

	return req, rr, eng
}
