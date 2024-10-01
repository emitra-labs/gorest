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
	gorest.Add(http.MethodGet, "/hello", sayHello)

	log.Fatal(gorest.Start())
}
