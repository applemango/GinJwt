package token

import (
	"errors"
	"fmt"

	"GinJwt/models"

	"github.com/dgrijalva/jwt-go"
)

func Login(username string, password string) (map[string]string, error) {
	user, err := models.GetUserFromUsername(username)
	//user, err := models.GetUserFromId(1)
	if err != nil {
		return nil, errors.New("user not found or invalid password")
	}
	fmt.Println(username, user.Username)
	fmt.Println(password, user.Password)

	if username == user.Username && password == user.Password {
		tokens, err := GenerateTokenPair(user.Id, user.Username)
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
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	fmt.Println(refresh_token, token, err)
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in

		user, error := models.GetUserFromId(int(claims["sub"].(float64)))
		if error != nil {
			return "", errors.New("user not found or invalid password")
		}

		if int(claims["sub"].(float64)) == user.Id {

			newTokenPair, err := GenerateTokenAccess(user.Id, user.Username)
			if err != nil {
				return "", err
			}
			return newTokenPair, nil
		}

		return "", errors.New("error")
	}

	return "", err

}
