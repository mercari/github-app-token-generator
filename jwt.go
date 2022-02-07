package main

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

// Function to sign a jwt via the private key
func issueJWTFromPEM(appID string, key *rsa.PrivateKey) string {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Add(-1 * time.Minute).Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:    appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return ss
}
