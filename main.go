package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"GinJwt/middleware"
	"GinJwt/models"
	"GinJwt/token"
)

func main() {
	type postUser struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	type postToken struct {
		RefreshToken string `form:"refresh_token"`
	}

	err := models.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	app := gin.Default()

	app.POST("/login", func(c *gin.Context) {
		var u postUser
		if c.ShouldBind(&u) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
			return
		}
		token, err := token.Login(u.Username, u.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	app.POST("/refresh", func(c *gin.Context) {
		var t postToken
		if c.ShouldBind(&t) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
			return
		}
		token, err := token.RefreshToken(t.RefreshToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	app.POST("/register", func(c *gin.Context) {
		var u postUser
		if c.ShouldBind(&u) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
			return
		}
		var user models.User
		user.Password = u.Password
		user.Username = u.Username

		_, err := models.InsertUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "created"})
	})

	app.POST("/logout", middleware.JwtRequired(true), func(c *gin.Context) {
		tokenString := strings.Split(c.Request.Header["Authorization"][0], " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, gin.H{"msg": err})
			return
		}

		models.InsertTokenBlockList(tokenString)
		c.JSON(http.StatusBadRequest, gin.H{"msg": "success"})
	})

	app.GET("/test", middleware.JwtRequired(false), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello!"})
	})
	app.GET("/testall", middleware.JwtRequired(true), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "hello!"})
	})

	app.Run(":4000")
}
