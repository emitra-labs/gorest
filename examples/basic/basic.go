package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/emitra-labs/common/types"
	"github.com/emitra-labs/gorest"
	// "github.com/emitra-labs/authz/permission"
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
		// Permission: permission.CreateTeam,
		RateLimit: gorest.RateLimit(3, 10*time.Second),
	})

	log.Fatal(gorest.Start())
}
