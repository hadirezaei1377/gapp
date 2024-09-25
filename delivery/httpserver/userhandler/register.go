package userhandler

import (
	dto "gapp/param"
	"gapp/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) userRegister(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if fieldErrors, err := h.userValidator.ValidateRegisterRequest(req); err != nil {
		msg, code := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.userSvc.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}
