package token

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func Login(username string, password string) (map[string]string, error) {

	if username == "test" && password == "test" {
		tokens, err := GenerateTokenPair()
		if err != nil {
			return nil, err
		}
		return tokens, nil
	}
	return nil, errors.New("user not found or invalid password")
}

func RefreshToken(refresh_token string) (string, error) {

	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(refresh_token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	fmt.Println(refresh_token, token, err)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in
		if int(claims["sub"].(float64)) == 1 {

			newTokenPair, err := GenerateTokenAccess()
			if err != nil {
				return "", err
			}
			return newTokenPair, nil
		}

		return "", errors.New("error")
	}

	return "", err

}
