package router

import (
	"github.com/labstack/echo/v4"

	v1 "goddd-boilerplate/app/handler/v1"
	"goddd-boilerplate/app/repository"
)

// versionOne :
func versionOne(e *echo.Echo, repo *repository.Repository) *echo.Echo {
	h := v1.New(repo)

	apiVersionOne := e.Group("/v1")
	apiVersionOne.GET("/", h.HealthCheck)
	apiVersionOne.POST("/data-migration", h.DataMigration) // data migration
	apiVersionOne.POST("/login", h.Login)
	apiVersionOne.POST("/register", h.Register)

	user := apiVersionOne.Group("", h.Authorization)
	user.GET("/profile", h.GetProfile)

	return e
}
