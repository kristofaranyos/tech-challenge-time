package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kristofaranyos/tech-challenge-time/repository"
	"log"
	"net/http"
	"strings"
)

type AuthMiddleware interface {
	Authenticate(c *gin.Context)
}

type authMiddleware struct {
	repo repository.TokenRepository
}

func NewAuthMiddleware(repo repository.TokenRepository) AuthMiddleware {
	return &authMiddleware{repo: repo}
}

func (am *authMiddleware) Authenticate(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
	token = strings.TrimPrefix(token, "Bearer ")

	valid, err := am.repo.IsTokenValid(c.Request.Context(), token)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if !valid {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("token", token)
}
