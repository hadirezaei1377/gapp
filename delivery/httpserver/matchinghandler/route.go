package matchinghandler

import (
	"gapp/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/matching")

	userGroup.POST("/add-to-waiting-list", h.addToWaitingList,
		middleware.Auth(h.authSvc, h.authConfig), middleware.UpsertPresence(h.presenceSvc))
}
