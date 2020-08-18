package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ParseTokenString(tokenString string, signingString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingString), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func ValidJWT(token *jwt.Token, requiredScopes map[string]interface{}) error {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		valid := claims.VerifyExpiresAt(time.Now().Unix(), true)
		if !valid {
			return errors.New("invalid token - expired")
		}

		scopes, ok := claims["scopes"].([]interface{})
		if !ok {
			return errors.New("invalid token - missing scopes")
		}

		for _, s := range scopes {
			sn := s.(string)
			if _, ok := requiredScopes[sn]; !ok {
				return errors.New("invalid token - not authorized scopes")
			}
		}

		return nil
	}

	return errors.New("invalid token")
}
