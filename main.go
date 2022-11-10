package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"GinJwt/models"
	token "GinJwt/token"
)

func main() {

	err := models.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	app := gin.Default()

	app.POST("/login", func(c *gin.Context) {
		type postUser struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}
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
		fmt.Printf("token: %v\n", token)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	app.POST("/refresh", func(c *gin.Context) {
		type postToken struct {
			RefreshToken string `form:"refresh_token"`
		}
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
		fmt.Printf("token: %v\n", token)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	app.POST("/register", func(c *gin.Context) {
		type postUser struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}
		var u postUser
		if c.ShouldBind(&u) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "error"})
			return
		}
		var user models.User
		user.Password = u.Password
		user.Username = u.Username

		result, err := models.InsertUser(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": err})
			return
		}
		fmt.Printf("result: %v\n", result)
		c.JSON(http.StatusOK, gin.H{"msg": "created"})
	})
	app.Run(":4000")
}
