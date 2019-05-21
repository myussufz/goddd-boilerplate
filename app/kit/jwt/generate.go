package jwt

import (
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// GenerateToken :
func GenerateToken(key []byte, claim jwt.Claims) (string, error) {
	enc, err := jose.NewSigner(
		jose.SigningKey{
			Algorithm: jose.HS256,
			Key:       key,
		},
		&jose.SignerOptions{
			ExtraHeaders: map[jose.HeaderKey]interface{}{
				"typ": "JWT",
			},
		},
	)
	if err != nil {
		return "", err
	}

	token, err := jwt.Signed(enc).Claims(claim).CompactSerialize()
	if err != nil {
		return "", err
	}

	return token, nil
}
