package v1

import (
	"encoding/json"

	"goddd-boilerplate/app/repository"
	"goddd-boilerplate/app/repository/kit"

	"github.com/labstack/echo/v4"
)

// Handler :
type Handler struct {
	repository *repository.Repository
}

// New :
func New(repo *repository.Repository) *Handler {
	return &Handler{
		repository: repo,
	}
}

// HealthCheck :
func (h Handler) HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]interface{}{
		"code": "SUCCESS",
	})
}

func paginate(c echo.Context, filterColumns ...string) (kit.Paginate, error) {
	paginator := kit.Paginate{}

	var i struct {
		Cursor  string `query:"cursor" json:"cursor"`
		Filters string `query:"filters" json:"filters"`
		Limit   int64  `query:"limit" json:"limit"`
	}

	if err := c.Bind(&i); err != nil {
		return paginator, err
	}

	filters := make(map[string]map[string]interface{})
	if err := json.Unmarshal([]byte(i.Filters), &filters); err != nil {
		return paginator, err
	}

	limit := int64(20)

	if i.Limit > 0 && i.Limit <= 200 {
		limit = i.Limit
	}

	paginator.Limit = limit
	paginator.Filters = filters
	paginator.FilterColumns = filterColumns

	return paginator, nil
}
