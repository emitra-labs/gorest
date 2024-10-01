package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/emitra-labs/common/types"
	"github.com/emitra-labs/gorest"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func createProduct(ctx context.Context, req *Product) (*types.BasicResponse, error) {
	return &types.BasicResponse{
		Message: fmt.Sprintf("New product created: %d", req.ID),
	}, nil
}

func main() {
	gorest.Add(http.MethodPost, "/product", createProduct)

	log.Fatal(gorest.Start())
}
