package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"mamlaka/cmd/web"
	_ "mamlaka/docs"
	"mamlaka/internal/app/payment"
	"mamlaka/internal/app/user"
	"net/http"
)

// @title Mamlaka API
// @version 0.1
// @description This is a sample server for an Echo API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// RegisterRoutes @host localhost:8080
// @BasePath /api/v1
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))

	e.GET("/assets/*", echo.WrapHandler(fileServer))
	// Swagger route
	e.GET("/docs/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", s.healthHandler)

	api := e.Group("/api/v1")
	{
		user.RegisterUserRoutes(api, s.logger, s.db.GetDB())
		payment.RegisterPaymentRoutes(api, s.logger, s.db.GetDB())
	}
	return e
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
