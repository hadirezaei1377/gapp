package matchinghandler

import (
	"gapp/param"
	"gapp/pkg/claim"
	"gapp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) addToWaitingList(c echo.Context) error {
	var req param.AddToWaitingListRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims := claim.GetClaimsFromEchoContext(c)
	req.UserID = claims.UserID

	if fieldErrors, err := h.matchingValidator.ValidateAddToWaitingListRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.matchingSvc.AddToWaitingList(req)
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
