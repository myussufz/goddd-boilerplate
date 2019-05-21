package transformer

import (
	"time"

	"goddd-boilerplate/app/model"
)

// User :
type User struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	CountryCode     string    `json:"country_code"`
	PhoneNumber     string    `json:"phone_number"`
	CreatedDateTime time.Time `json:"createdAt"`
	UpdatedDateTime time.Time `json:"updatedAt"`
}

// ToUser :
func ToUser(user *model.User) *User {
	u := new(User)
	u.ID = user.ID.Hex()
	u.Name = user.Name
	u.CountryCode = user.CountryCode
	u.PhoneNumber = user.PhoneNumber
	u.CreatedDateTime = user.CreatedDateTime
	u.UpdatedDateTime = user.UpdatedDateTime

	return u
}
