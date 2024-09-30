package backofficeuserhandler

import (
	"gapp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) listUsers(c echo.Context) error {
	list, err := h.backofficeUserSvc.ListAllUsers()
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}
	return c.JSON(http.StatusOK, echo.Map{
		"data": list,
	})
}
