package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kristofaranyos/tech-challenge-time/controller/dto"
	"github.com/kristofaranyos/tech-challenge-time/repository"
	"github.com/kristofaranyos/tech-challenge-time/util"
	"log"
	"net/http"
	"time"
)

type SessionController interface {
	Start(c *gin.Context)
	Stop(c *gin.Context)
	UpdateName(c *gin.Context)
	List(c *gin.Context)
}

type sessionController struct {
	repo repository.SessionRepository
}

func NewSessionController(repo repository.SessionRepository) SessionController {
	return &sessionController{repo: repo}
}

func (sc *sessionController) Start(c *gin.Context) {
	// Token definitely exists because this is an authenticated endpoint
	tokenIntf, _ := c.Get("token")

	sessionId, err := sc.repo.StartSession(c.Request.Context(), tokenIntf.(string), util.GetUTCDateTime())
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessionId": sessionId})
}

func (sc *sessionController) Stop(c *gin.Context) {
	var sessionStruct dto.StopRequest
	if err := c.ShouldBindJSON(&sessionStruct); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	session, err := sc.repo.GetSession(c.Request.Context(), sessionStruct.SessionId)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Token definitely exists because this is an authenticated endpoint
	tokenIntf, _ := c.Get("token")

	// Users can only stop their own sessions
	if tokenIntf.(string) != session.User {
		c.Status(http.StatusUnauthorized)
		return
	}

	// Check if already stopped
	if session.Stopped != "1000-01-01 00:00:00" {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := sc.repo.StopSession(c.Request.Context(), sessionStruct.SessionId, util.GetUTCDateTime()); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (sc *sessionController) UpdateName(c *gin.Context) {
	var nameUpdateStruct dto.NameUpdateRequest
	if err := c.ShouldBindJSON(&nameUpdateStruct); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	session, err := sc.repo.GetSession(c.Request.Context(), nameUpdateStruct.SessionId)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Token definitely exists because this is an authenticated endpoint
	tokenIntf, _ := c.Get("token")

	// Users can only stop their own sessions
	if tokenIntf.(string) != session.User {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := sc.repo.SetName(c.Request.Context(), nameUpdateStruct.SessionId, nameUpdateStruct.Name); err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

func (sc *sessionController) List(c *gin.Context) {
	// Token definitely exists because this is an authenticated endpoint
	tokenIntf, _ := c.Get("token")

	var from, to string
	switch c.Param("duration") {
	case "": // If empty, no date restriction is set
	case "day":
		from, to = util.FormatDateTime(time.Now().UTC().Add(-24*time.Hour)), util.GetUTCDateTime()
	case "week":
		from, to = util.FormatDateTime(time.Now().UTC().Add(-24*7*time.Hour)), util.GetUTCDateTime()
	case "month":
		from, to = util.FormatDateTime(time.Now().UTC().Add(-24*7*4*time.Hour)), util.GetUTCDateTime()
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid listing interval"})
		return
	}

	list, err := sc.repo.ListSessions(c.Request.Context(), tokenIntf.(string), from, to)
	if err != nil {
		log.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	response := make([]dto.ListResponse, len(list))
	for i := 0; i < len(list); i++ {
		response[i].Fill(list[i])
	}

	c.JSON(http.StatusOK, response)
}
