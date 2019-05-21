package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// DataMigration :
func (h Handler) DataMigration(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}
