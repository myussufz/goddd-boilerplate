package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"goddd-boilerplate/app/kit/validator"
	"goddd-boilerplate/app/response"
	"goddd-boilerplate/app/response/errcode"
	"goddd-boilerplate/app/router"
)

// Start :
func Start(port string) {
	e := echo.New()
	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return fmt.Sprintf("%d%d", time.Now().UnixNano(), rand.Intn(100000))
			},
		}),
		middleware.Secure(),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:9000", "http://localhost:3000"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		}),
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			bb := new(bytes.Buffer)
			bb.WriteString(`{`)
			bb.WriteString(`"time":"` + time.Now().Format(time.RFC3339Nano) + `",`)
			bb.WriteString(`"id":"` + c.Response().Header().Get("X-Request-Id") + `",`)

			request := new(bytes.Buffer)
			if err := json.Compact(request, reqBody); err == nil {
				bb.WriteString(`"request":`)
				bb.WriteString(request.String() + ",")
			}
			bb.WriteString(`"response":`)

			response := string(bytes.TrimSpace(resBody))
			if response == "" {
				bb.WriteString("{}")
			} else {
				bb.WriteString(response)
			}
			bb.WriteString(`}`)
			bb.WriteString("\n")
			os.Stdout.Write(bb.Bytes())
		}),
	)

	// register custom validator
	e.Validator = validator.New()

	// register new error handler
	e.HTTPErrorHandler = customErrorHandler

	// init new router
	e = router.New(e)

	e.Logger.Fatal(e.Start(":" + port))
}

// customErrorHandler :
func customErrorHandler(err error, c echo.Context) {
	code := strings.TrimSpace(strings.Replace(err.Error(), "code=", "", -1))

	switch code[:3] {
	case "404", "405":
		c.JSON(http.StatusNotFound,
			response.Exception{
				Code:  errcode.APIEndpointNotExist,
				Error: errors.WithStack(err),
			})

	default:
		c.JSON(http.StatusInternalServerError,
			response.Exception{
				Code:  errcode.InvalidRequest,
				Error: errors.WithStack(err),
			})
	}
}
