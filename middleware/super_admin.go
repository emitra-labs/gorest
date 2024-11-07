package middleware

import (
	"github.com/emitra-labs/common/constant"
	"github.com/emitra-labs/common/errors"
	"github.com/labstack/echo/v4"
)

func SuperAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			superAdmin, _ := ctx.Value(constant.SuperAdmin).(bool)

			if !superAdmin {
				return errors.PermissionDenied()
			}

			return next(c)
		}
	}
}
