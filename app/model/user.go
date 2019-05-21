package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User :
type User struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string             `bson:"name"`
	CountryCode     string             `bson:"country_code"`
	PhoneNumber     string             `bson:"phone_number"`
	PasswordHash    string             `bson:"password_hash"`
	PasswordSalt    string             `bson:"password_salt"`
	LastSignedAt    time.Time          `json:"last_signed_at"`
	CreatedDateTime time.Time          `bson:"created_at"`
	UpdatedDateTime time.Time          `bson:"updated_at"`
}
