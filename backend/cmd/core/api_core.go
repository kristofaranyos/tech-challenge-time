package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kristofaranyos/tech-challenge-time/controller"
	"github.com/kristofaranyos/tech-challenge-time/middleware"
	"github.com/kristofaranyos/tech-challenge-time/repository"
	"github.com/magiconair/properties"
	"time"
)

type apiCore struct {
	cfg    *ApiConfig
	router *gin.Engine
	db     *sqlx.DB
}

func NewApiCore() (*apiCore, error) {
	core := &apiCore{}
	core.cfg = &ApiConfig{}

	p := properties.MustLoadFile("app.properties", properties.UTF8)
	if err := p.Decode(core.cfg); err != nil {
		return nil, fmt.Errorf("Failed to read properties file: %v", err)
	}

	if err := core.initDB(); err != nil {
		return nil, err
	}

	core.initRoutes()

	return core, nil
}

func (c *apiCore) Run() error {
	return c.router.Run(":" + c.cfg.Port)
}

func (c *apiCore) Close() error {
	return c.db.Close()
}

func (c *apiCore) initRoutes() {
	gin.SetMode(gin.ReleaseMode)
	c.router = gin.Default()
	api := c.router.Group("/api/v1")

	tokenRepository := repository.NewTokenRepository(c.db)

	{
		cont := controller.NewAuthController(tokenRepository)
		api.GET("/auth/newtoken", cont.NewToken)
	}

	{
		session := api.Group("/session")

		auth := middleware.NewAuthMiddleware(tokenRepository)
		session.Use(auth.Authenticate)

		cont := controller.NewSessionController(repository.NewSessionRepository(c.db))
		session.POST("/start", cont.Start)
		session.POST("/stop", cont.Stop)
		session.PUT("/name", cont.UpdateName)
		session.GET("/list", cont.List)
		session.GET("/list/:duration", cont.List)
	}
}

func (c *apiCore) initDB() error {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:@%s(%s)/%s", c.cfg.DbUSer, c.cfg.DbPass, c.cfg.DbHost, c.cfg.DbName))
	if err != nil {
		return fmt.Errorf("Failed to connect to MySQL: %v", err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)

	c.db = db

	return nil
}
