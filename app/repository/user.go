package repository

import (
	"goddd-boilerplate/app/model"
	"goddd-boilerplate/app/repository/kit"
)

// User :
type User interface {
	Create(user *model.User) error
	FindByPhoneNumber(countryCode, phoneNumber string) (*model.User, error)
	UniqueIndex(key, value string) error
	Paginate(p kit.Paginate) ([]*model.User, string, error)
	FindByID(id string) (*model.User, error)
}
