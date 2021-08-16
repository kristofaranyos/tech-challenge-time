package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kristofaranyos/tech-challenge-time/repository"
	"github.com/kristofaranyos/tech-challenge-time/util"
	"log"
	"net/http"
)

type AuthController interface {
	NewToken(c *gin.Context)
}

type authController struct {
	repo repository.TokenRepository
}

func NewAuthController(repo repository.TokenRepository) AuthController {
	return &authController{repo: repo}
}

func (ac *authController) NewToken(c *gin.Context) {
	token := uuid.New().String()

	if err := ac.repo.InsertToken(c.Request.Context(), token, util.GetUTCDateTime()); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
