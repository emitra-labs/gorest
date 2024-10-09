package main

import (
	"context"
	"log"
	"net/http"

	"github.com/emitra-labs/common/types"
	"github.com/emitra-labs/gorest"
)

func sayHello(ctx context.Context, req *types.Empty) (*types.BasicResponse, error) {
	return &types.BasicResponse{
		Message: "Hello, World!",
	}, nil
}

func main() {
	gorest.Add(http.MethodGet, "/hello", sayHello, gorest.RouteConfig{
		Summary:     "Say hello",
		Description: `Say hello to the world`,
		Tags:        []string{"Greeting"},
	})

	log.Fatal(gorest.Start())
}
