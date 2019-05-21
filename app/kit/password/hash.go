package password

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Create :
func Create(password, salt, pepper string) (string, error) {
	passwordPepper := fmt.Sprintf("%s%s%s", password, salt, pepper)
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(passwordPepper), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hashPassword), nil
}

// Compare : Check whether password is match or not
func Compare(password, salt, pepper, hashedPassword string) bool {
	hashPassword, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return false
	}
	p := []byte(fmt.Sprintf("%s%s%s", password, salt, pepper))
	if err := bcrypt.CompareHashAndPassword(hashPassword, p); err != nil {
		return false
	}
	return true
}
