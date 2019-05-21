package v1

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	zxcvbn "github.com/nbutton23/zxcvbn-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	jose "gopkg.in/square/go-jose.v2/jwt"

	"goddd-boilerplate/app/config"
	"goddd-boilerplate/app/kit/jwt"
	"goddd-boilerplate/app/kit/password"
	"goddd-boilerplate/app/kit/random"
	"goddd-boilerplate/app/model"
	"goddd-boilerplate/app/response"
	"goddd-boilerplate/app/response/errcode"
)

const (
	tokenExpiresInSecond = 86400
)

// Register :
func (h Handler) Register(c echo.Context) error {
	var i struct {
		CountryCode string `json:"countryCode" validate:"required,numeric,min=2,max=2"`
		PhoneNumber string `json:"phoneNumber" validate:"required,numeric,min=7,max=12"`
		Name        string `json:"name" validate:"required"`
		Password    string `json:"password" validate:"required"`
	}

	if err := c.Bind(&i); err != nil {
		return c.JSON(http.StatusBadRequest,
			response.Exception{
				Code:  errcode.InvalidRequest,
				Error: errors.WithStack(err),
			})
	}

	i.CountryCode = strings.TrimSpace(i.CountryCode)
	i.PhoneNumber = strings.TrimSpace(strings.TrimLeft(i.PhoneNumber, "0"))
	i.Password = strings.TrimSpace(i.Password)

	if err := c.Validate(&i); err != nil {
		return c.JSON(http.StatusUnprocessableEntity,
			response.Exception{
				Code:   errcode.ValidationError,
				Detail: err.Error(),
				Error:  errors.WithStack(err),
			})
	}

	strength := zxcvbn.PasswordStrength(i.Password, nil).Score
	if strength < 2 {
		return c.JSON(http.StatusExpectationFailed,
			response.Exception{
				Code: errcode.PasswordNotStrength,
			})
	}

	salt := random.String(10)
	pepper := config.SecretKey

	// take 80 ms to do password hash
	passwordHash, err := password.Create(i.Password, salt, pepper)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			response.Exception{
				Code:  errcode.SystemError,
				Error: errors.WithStack(err),
			})
	}

	if _, err := h.repository.User.FindByPhoneNumber(i.CountryCode, i.PhoneNumber); err == nil {
		return c.JSON(http.StatusConflict,
			response.Exception{
				Code: errcode.PhoneNumberAlreadyExist,
			})
	}

	user := new(model.User)
	user.ID = primitive.NewObjectID()
	user.Name = i.Name
	user.CountryCode = i.CountryCode
	user.PhoneNumber = i.PhoneNumber
	user.PasswordHash = passwordHash
	user.PasswordSalt = salt
	user.LastSignedAt = time.Now().UTC()
	user.CreatedDateTime = time.Now().UTC()
	user.UpdatedDateTime = time.Now().UTC()

	if err := h.repository.User.Create(user); err != nil {
		return c.JSON(http.StatusInternalServerError,
			response.Exception{
				Code:  errcode.SystemError,
				Error: errors.WithStack(err),
			})
	}

	expiredDateTime := time.Now().UTC().Add(time.Duration(tokenExpiresInSecond) * time.Second)

	accessToken, err := jwt.GenerateToken(config.JWTKey, jose.Claims{
		Issuer:    config.SystemPath,
		Subject:   user.ID.Hex(),
		Audience:  []string{user.ID.Hex()},
		Expiry:    jose.NewNumericDate(expiredDateTime),
		NotBefore: jose.NewNumericDate(time.Now().UTC()),
		IssuedAt:  jose.NewNumericDate(time.Now().UTC()),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			response.Exception{
				Code:  errcode.InvalidRequest,
				Error: errors.WithStack(err),
			})
	}

	return c.JSON(http.StatusAccepted, response.Item{
		Item: map[string]interface{}{
			"accessToken": accessToken,
			"expiresIn":   tokenExpiresInSecond,
			"expiredAt":   expiredDateTime,
		},
	})
}

// Login :
func (h Handler) Login(c echo.Context) error {
	var i struct {
		CountryCode string `json:"countryCode" validate:"required,numeric,min=2,max=2"`
		PhoneNumber string `json:"phoneNumber" validate:"required,numeric,min=7,max=12"`
		Password    string `json:"password" validate:"required"`
	}

	if err := c.Bind(&i); err != nil {
		return c.JSON(http.StatusBadRequest,
			response.Exception{
				Code:  errcode.InvalidRequest,
				Error: errors.WithStack(err),
			})
	}

	i.CountryCode = strings.TrimSpace(i.CountryCode)
	i.PhoneNumber = strings.TrimSpace(strings.TrimLeft(i.PhoneNumber, "0"))
	i.Password = strings.TrimSpace(i.Password)

	if err := c.Validate(&i); err != nil {
		return c.JSON(http.StatusUnprocessableEntity,
			response.Exception{
				Code:   errcode.ValidationError,
				Detail: err.Error(),
				Error:  errors.WithStack(err),
			})
	}

	user, err := h.repository.User.FindByPhoneNumber(i.CountryCode, i.PhoneNumber)
	if err != nil {
		return c.JSON(http.StatusForbidden,
			response.Exception{
				Code: errcode.UserFailedAuthentication,
			})
	}

	salt := user.PasswordSalt
	pepper := config.SecretKey

	if isExist := password.Compare(i.Password, salt, pepper, user.PasswordHash); !isExist {
		return c.JSON(http.StatusForbidden,
			response.Exception{
				Code: errcode.UserFailedAuthentication,
			})
	}

	expiredDateTime := time.Now().UTC().Add(time.Duration(tokenExpiresInSecond) * time.Second)

	accessToken, err := jwt.GenerateToken(config.JWTKey, jose.Claims{
		Issuer:    config.SystemPath,
		Subject:   user.ID.Hex(),
		Audience:  []string{user.ID.Hex()},
		Expiry:    jose.NewNumericDate(expiredDateTime),
		NotBefore: jose.NewNumericDate(time.Now().UTC()),
		IssuedAt:  jose.NewNumericDate(time.Now().UTC()),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			response.Exception{
				Code:  errcode.InvalidRequest,
				Error: errors.WithStack(err),
			})
	}

	return c.JSON(http.StatusAccepted, response.Item{
		Item: map[string]interface{}{
			"accessToken": accessToken,
			"expiresIn":   tokenExpiresInSecond,
		},
	})
}

// Authorization :
func (h Handler) Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenArr := strings.Split(strings.TrimSpace(c.Request().Header.Get("Authorization")), " ")
		if len(tokenArr) != 2 {
			return c.JSON(http.StatusUnauthorized,
				response.Exception{
					Code: errcode.TokenInvalid,
				})
		}

		switch strings.TrimSpace(tokenArr[0]) {
		case "Bearer":
		default:
			return c.JSON(http.StatusUnauthorized,
				response.Exception{
					Code: errcode.TokenInvalid,
				})
		}

		token := strings.TrimSpace(tokenArr[1])
		claim, isValid := jwt.ValidateToken(config.JWTKey, token)
		if !isValid {
			return c.JSON(http.StatusUnauthorized,
				response.Exception{
					Code: errcode.TokenInvalid,
				})
		}

		if err := claim.Validate(
			jose.Expected{
				Issuer: config.SystemPath,
			},
		); err != nil {
			return c.JSON(http.StatusUnauthorized,
				response.Exception{
					Code: errcode.TokenInvalid,
				})
		}

		user, err := h.repository.User.FindByID(claim.Subject)
		if err != nil {
			return c.JSON(http.StatusUnauthorized,
				response.Exception{
					Code: errcode.TokenInvalid,
				})
		}

		c.Set(model.CollectionUser, *user)

		return next(c)
	}
}
