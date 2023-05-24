package token

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type JWTTokenVerifier struct {
	PublicKey *rsa.PublicKey
}

// Verify 验证token
func (v *JWTTokenVerifier) Verify(token string) (string, error) {
	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return v.PublicKey, nil
		})
	if err != nil {
		return "", fmt.Errorf("cannot parse token:%v", err)
	}
	if !t.Valid {
		return "", fmt.Errorf("token not valid")
	}

	claims, ok := t.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("token claim is not StandardClaim")
	}
	err = claims.Valid()
	if err != nil {
		return "", fmt.Errorf("claim not valid:%v", err)
	}

	return claims.Subject, nil
}
