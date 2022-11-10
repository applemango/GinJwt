package middleware

import (
	"GinJwt/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JwtRequired(refresh bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		fmt.Println(models.TokenInTokenBlockList(tokenString))
		if err != nil || !token.Valid || (!refresh && token.Claims.(jwt.MapClaims)["refresh"] == true) || models.TokenInTokenBlockList(tokenString) {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
			c.Abort()
		}
	}
}
