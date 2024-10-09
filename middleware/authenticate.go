package middleware

import (
	"context"
	"strings"

	"github.com/emitra-labs/authn"
	"github.com/emitra-labs/common/constant"
	"github.com/emitra-labs/common/errors"
	"github.com/labstack/echo/v4"
)

func Authenticate() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")

			if accessToken == "" {
				return errors.Unauthenticated("No access token provided")
			}

			claims, err := authn.VerifyToken(accessToken)
			if err != nil {
				return err
			}

			ctx := context.WithValue(c.Request().Context(), constant.UserID, claims.Subject)
			ctx = context.WithValue(ctx, constant.SessionID, claims.SessionID)
			ctx = context.WithValue(ctx, constant.SuperAdmin, claims.SuperAdmin)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
