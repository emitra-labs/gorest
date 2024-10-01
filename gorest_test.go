package gorest_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/emitra-labs/common/types"
	"github.com/emitra-labs/gorest"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

var handler http.Handler

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

func TestMain(m *testing.M) {
	gorest.Add(http.MethodGet, "/hello", sayHello)

	gorest.Add(http.MethodPost, "/product", createProduct)

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
