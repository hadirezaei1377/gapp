package userhandler

import "github.com/labstack/echo/v4"

func (h Handler) SetUserRoutes(e *echo.Echo) {
	userGroup := e.Group("/users")

	userGroup.GET("/profile", h.userProfile)
	userGroup.POST("/login", h.userLogin)
	userGroup.POST("/register", h.userRegister)
}
