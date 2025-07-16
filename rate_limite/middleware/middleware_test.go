package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dlcdev1/pos1/rate_limite/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockLimiter struct {
	allow   bool
	err     error
	lastIP  string
	lastTok string
}

func (m *mockLimiter) Allow(token, ip string) (bool, error) {
	m.lastIP = ip
	m.lastTok = token
	return m.allow, m.err
}

func TestRateLimiterMiddleware_Allowed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ml := &mockLimiter{allow: true, err: nil}
	r := gin.New()

	r.Use(middleware.RateLimiterMiddleware(ml))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "allowed"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "token123")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message":"allowed"}`, w.Body.String())
	assert.Equal(t, "token123", ml.lastTok)
}

func TestRateLimiterMiddleware_Blocked(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ml := &mockLimiter{allow: false, err: nil}
	r := gin.New()
	r.Use(middleware.RateLimiterMiddleware(ml))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not be visible"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "token123")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Contains(t, w.Body.String(), "you have reached the maximum number of requests")
	assert.Equal(t, "token123", ml.lastTok)
}

func TestRateLimiterMiddleware_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ml := &mockLimiter{allow: false, err: errors.New("redis failure")}
	r := gin.New()

	r.Use(middleware.RateLimiterMiddleware(ml))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "should not be visible"})
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "token123")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.JSONEq(t, `{"error":"failed to validate rate limit"}`, w.Body.String())
	assert.Equal(t, "token123", ml.lastTok)
}
