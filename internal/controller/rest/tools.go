package rest

import (
	"github.com/go-chi/jwtauth/v5"
)

const _secret = "<jwt-secret>"

var TokenAuth *jwtauth.JWTAuth

func init() {
	TokenAuth = jwtauth.New("HS256", []byte(_secret), nil)
}

func MakeToken(userID string) (string, error) {
	/*
		Create a JWT token.
	*/
	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{"user_id": userID})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
