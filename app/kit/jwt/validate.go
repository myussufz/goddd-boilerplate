package jwt

import (
	"gopkg.in/square/go-jose.v2/jwt"
)

// ValidateToken :
func ValidateToken(key []byte, token string) (*jwt.Claims, bool) {
	enc, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, false
	}

	claim := new(jwt.Claims)
	if err := enc.Claims(key, claim); err != nil {
		return nil, false
	}

	return claim, true
}
