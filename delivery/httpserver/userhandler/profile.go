package userhandler

import (
	"gapp/config"
	"gapp/param"
	"gapp/pkg/httpmsg"
	"gapp/service/authservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func getClaims(c echo.Context) *authservice.Claims {
	return c.Get(config.AuthMiddlewareContextKey).(*authservice.Claims)
}

func (h Handler) userProfile(c echo.Context) error {
	claims := getClaims(c)

	resp, err := h.userSvc.Profile(param.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
