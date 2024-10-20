package middleware

import (
	"gapp/entity"
	"gapp/pkg/claim"
	"gapp/pkg/errmsg"
	"gapp/service/authorizationservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccessCheck(service authorizationservice.Service,
	permissions ...entity.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)
			if err != nil {
				// TODO - log unexpected error
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorMsgSomethingWentWrong,
				})
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorMsgUserNotAllowed,
				})
			}

			return next(c)
		}
	}
}
