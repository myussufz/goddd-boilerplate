package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"goddd-boilerplate/app/model"
	"goddd-boilerplate/app/response"
	"goddd-boilerplate/app/response/errcode"
	"goddd-boilerplate/app/response/transformer"
)

// GetProfile :
func (h Handler) GetProfile(c echo.Context) error {
	user := c.Get(model.CollectionUser).(model.User)

	return c.JSON(http.StatusAccepted, response.Item{
		Item: transformer.ToUser(&user),
	})
}

// GetUserProfile : Paginate example
func (h Handler) GetUserProfile(c echo.Context) error {
	paginator, err := paginate(c, "name")
	if err != nil {
		return c.JSON(http.StatusBadRequest,
			response.Exception{
				Code:  errcode.ValidationError,
				Error: errors.WithStack(err),
			})
	}

	users, cursor, err := h.repository.User.Paginate(paginator)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			response.Exception{
				Code:  errcode.SystemError,
				Error: errors.WithStack(err),
			})
	}

	return c.JSON(http.StatusOK,
		response.Items{
			Items: funk.Map(users, func(user *model.User) *transformer.User {
				return transformer.ToUser(user)
			}),
			Cursor: cursor,
		})
}
