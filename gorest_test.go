package gorest_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/emitra-labs/authn"
	"github.com/emitra-labs/common/constant"
	"github.com/emitra-labs/common/types"
	"github.com/emitra-labs/gorest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var handler http.Handler
var accessToken string

func sayHello(ctx context.Context, req *types.Empty) (*types.BasicResponse, error) {
	return &types.BasicResponse{
		Message: "Hello, World!",
	}, nil
}

func createProduct(ctx context.Context, req *Product) (*types.BasicResponse, error) {
	return &types.BasicResponse{
		Message: fmt.Sprintf("New product created: %d", req.ID),
	}, nil
}

func restricted(ctx context.Context, req *types.Empty) (*types.BasicResponse, error) {
	userID, _ := ctx.Value(constant.UserID).(string)

	return &types.BasicResponse{
		Message: fmt.Sprintf("UserID: %s", userID),
	}, nil
}

func TestMain(m *testing.M) {
	accessToken, _ = authn.GenerateToken(authn.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "john",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		},
		SessionID:  "123",
		SuperAdmin: true,
	})

	gorest.Add(http.MethodGet, "/hello", sayHello)

	gorest.Add(http.MethodPost, "/product", createProduct)

	gorest.Add(http.MethodGet, "/restricted", restricted, gorest.RouteConfig{
		Authenticate: true,
	})

	handler = gorest.GetHandler()

	os.Exit(m.Run())
}

func TestHello(t *testing.T) {
	apitest.New().
		Handler(handler).
		Get("/hello").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.message", "Hello, World!")).
		End()
}

func TestCreateProduct(t *testing.T) {
	apitest.New().
		Handler(handler).
		Post("/product").
		JSON(`{"id": 73, "name": "Product 73", "price": 35}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.message", "New product created: 73")).
		End()
}

func TestNotFound(t *testing.T) {
	apitest.New().
		Handler(handler).
		Get("/not-found").
		Expect(t).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.error", "Not Found")).
		End()
}

func TestRestricted_Success(t *testing.T) {
	apitest.New().
		Handler(handler).
		Get("/restricted").
		Header("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Equal("$.message", "UserID: john")).
		End()
}

func TestRestricted_Unauthenticated(t *testing.T) {
	apitest.New().
		Handler(handler).
		Get("/restricted").
		Expect(t).
		Status(http.StatusUnauthorized).
		End()
}
