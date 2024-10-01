# gorest

Quickly build RESTful APIs in Go with auto-generated OpenAPI 3.1 documentation.

## Features

- Auto-parse body and parameters using struct field tags, e.g. `json` and `query`.
- Auto-generate OpenAPI 3.1 documentation.

## Quickstart

Make sure the environment variables are set (refer to `.env.example`).

Here's a simple code that demonstrates how to use it:

```go
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
```

Then, run the program:

```bash
go run server.go
```

Browse to http://localhost:3000/docs and the OpenAPI documentation should be displayed.
