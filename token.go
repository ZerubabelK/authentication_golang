package main

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(userID string) (string, error) {
    token := jwt.New(jwt.SigningMethodHS256)

    claims := token.Claims.(jwt.MapClaims)
    claims["user_id"] = userID
    claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()

    tokenString, err := token.SignedString([]byte("sfyTVOBxGNRslpohYhUQGTUlti09gHmFmDfBivvua6Q="))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("sfyTVOBxGNRslpohYhUQGTUlti09gHmFmDfBivvua6Q="), nil
	})
	log.Println(claims)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	return claims, nil
}