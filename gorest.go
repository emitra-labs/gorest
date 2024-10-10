package gorest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/emitra-labs/common/errors"
	"github.com/emitra-labs/common/validator"
	"github.com/emitra-labs/gorest/middleware"
	"github.com/labstack/echo/v4"
	"github.com/sethvargo/go-envconfig"
	"github.com/swaggest/openapi-go/openapi31"
)

type Server struct {
	Echo    *echo.Echo
	OpenAPI OpenAPI
}

type Config struct {
	Port      int `env:"GOREST_PORT, default=3000"`
	Info      Info
	ServerURL string `env:"GOREST_SERVER_URL, default=http://localhost:3000"`
}

type Info struct {
	Title       string `env:"GOREST_INFO_TITLE, default=My API" json:"title"`
	Description string `env:"GOREST_INFO_DESCRIPTION, default=This is a sample RESTful API server." json:"description,omitempty"`
	Version     string `env:"GOREST_INFO_VERSION, default=1.0.0" json:"version"`
}

type OpenAPI struct {
	Reflector *openapi31.Reflector
}

var config *Config
var server *Server

func init() {
	config = new(Config)

	// Load config from environment variables
	err := envconfig.Process(context.Background(), config)
	if err != nil {
		panic(err)
	}

	// Validate config
	err = validator.Validate(config)
	if err != nil {
		panic(err)
	}

	// Create server
	server = &Server{
		Echo: echo.New(),
		OpenAPI: OpenAPI{
			Reflector: &openapi31.Reflector{
				Spec: &openapi31.Spec{
					Openapi: "3.1.0",
					Info: openapi31.Info{
						Title:       config.Info.Title,
						Description: &config.Info.Description,
						Version:     config.Info.Version,
					},
					Servers: []openapi31.Server{{URL: config.ServerURL}},
				},
			},
		},
	}

	server.OpenAPI.Reflector.Spec.SetHTTPBearerTokenSecurity("Bearer token", "JWT", "")

	server.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		message := "Internal Server Error"

		if e, ok := err.(*errors.Error); ok {
			code = e.GetHTTPStatus()
			message = e.Error()
		} else if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			message = he.Message.(string)
		}

		c.JSON(code, map[string]string{
			"error": message,
		})
	}
}

func GetServer() *Server {
	return server
}

const swaggerHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>%s</title>
    <link
      rel="stylesheet"
      href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css"
    />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script
      src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"
      crossorigin
    ></script>
    <script>
      window.onload = () => {
        window.ui = SwaggerUIBundle({
          url: "./openapi.json",
          dom_id: "#swagger-ui",
        });
      };
    </script>
  </body>
</html>`

func Start() error {
	specBytes, _ := server.OpenAPI.Reflector.Spec.MarshalJSON()
	spec := string(specBytes)

	// Remove all package name prefixes
	spec = strings.ReplaceAll(spec, "schemas/Model", "schemas/")
	spec = strings.ReplaceAll(spec, "schemas/Types", "schemas/")
	spec = strings.ReplaceAll(spec, "\"Model", "\"")
	spec = strings.ReplaceAll(spec, "\"Types", "\"")

	var specJSON any
	json.Unmarshal([]byte(spec), &specJSON)

	server.Echo.GET("/docs", func(c echo.Context) error {
		return c.HTML(http.StatusOK, fmt.Sprintf(swaggerHTML, config.Info.Title))
	})

	server.Echo.GET("/openapi.json", func(c echo.Context) error {
		return c.JSON(http.StatusOK, specJSON)
	})

	server.Echo.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	return server.Echo.Start(fmt.Sprintf(":%d", config.Port))
}

func Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return server.Echo.Shutdown(ctx)
}

type RouteConfig struct {
	Summary      string
	Description  string
	Tags         []string
	Authenticate bool
	SuperAdmin   bool
	Permission   string
}

func Add[I, O any](
	method string,
	path string,
	f func(context.Context, *I) (*O, error),
	configs ...RouteConfig,
) {
	in := new(I)
	out := new(O)
	config := RouteConfig{}
	middlewares := []echo.MiddlewareFunc{}

	if len(configs) > 0 {
		config = configs[0]
	}

	op, _ := server.OpenAPI.Reflector.NewOperationContext(method, formatOpenAPIPath(path))
	op.SetSummary(config.Summary)
	op.SetDescription(config.Description)
	op.SetTags(config.Tags...)
	op.AddReqStructure(in)
	op.AddRespStructure(out)

	if config.Authenticate || config.SuperAdmin || config.Permission != "" {
		middlewares = append(middlewares, middleware.Authenticate())
		op.AddSecurity("Bearer token")
	}

	if config.SuperAdmin {
		middlewares = append(middlewares, middleware.SuperAdmin())
	}

	if err := server.OpenAPI.Reflector.AddOperation(op); err != nil {
		panic(err)
	}

	server.Echo.Add(method, path, func(c echo.Context) error {
		in := new(I)

		err := c.Bind(in)
		if err != nil {
			return err
		}

		res, err := f(c.Request().Context(), in)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	}, middlewares...)
}

func GetHandler() *echo.Echo {
	return server.Echo
}

func formatOpenAPIPath(path string) string {
	re := regexp.MustCompile(`:(\w+)`)
	return re.ReplaceAllString(path, `{$1}`)
}
