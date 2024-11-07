package middleware

import (
	"context"
	"crypto"
	"strings"

	"github.com/emitra-labs/common/constant"
	"github.com/emitra-labs/common/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	SessionID  string `json:"sid,omitempty"`
	SuperAdmin bool   `json:"adm,omitempty"`
}

func Authenticate(publicKey crypto.PublicKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			accessToken := c.Request().Header.Get("Authorization")
			accessToken = strings.TrimPrefix(accessToken, "Bearer ")

			if accessToken == "" {
				return errors.Unauthenticated("No access token provided")
			}

			parsed, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err != nil {
				return errors.Unauthenticated()
			}

			if claims, ok := parsed.Claims.(*Claims); ok {
				ctx := context.WithValue(c.Request().Context(), constant.UserID, claims.Subject)
				ctx = context.WithValue(ctx, constant.SessionID, claims.SessionID)
				ctx = context.WithValue(ctx, constant.SuperAdmin, claims.SuperAdmin)

				c.SetRequest(c.Request().WithContext(ctx))

				return next(c)
			}

			return errors.Unauthenticated()
		}
	}
}
