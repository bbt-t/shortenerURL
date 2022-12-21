package app

import "github.com/go-chi/jwtauth/v5"

const _secret = "<jwt-secret>"

var _tokenAuth *jwtauth.JWTAuth

func init() {
	_tokenAuth = jwtauth.New("HS256", []byte(_secret), nil)
}

func makeToken(name string) (string, error) {
	/*
		Create a JWT token.
	*/
	_, tokenString, err := _tokenAuth.Encode(map[string]interface{}{"username": name})
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
