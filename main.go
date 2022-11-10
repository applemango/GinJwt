package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	token "GinJwt/token"
)

func main() {

	app := gin.Default()

	/*app.GET("/url/get", func(c *gin.Context) {
		id := c.Query("id")
		intId, err_change := strconv.Atoi(id)
		data, err_get := models.GetUrl(intId)
		if id == "" || err_change != nil || err_get != nil {
			c.JSON(http.StatusOK, gin.H{
				"id":  -1,
				"url": "https://example.com",
			})
		}
		c.JSON(http.StatusOK, gin.H{
			"id":  data.Id,
			"url": data.Url,
		})
	})*/
	/*app.POST("/url/create", func(c *gin.Context) {
		var u postUrl
		if c.ShouldBind(&u) != nil {
			c.JSON(http.StatusOK, gin.H{
				"id":  -1,
				"url": "https://example.com",
			})
		}
		add(u.Url)
		url, _ := models.GetLastUrl()
		c.JSON(http.StatusOK, gin.H{
			"id":  url.Id,
			"url": url.Url,
		})
	})*/

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

	app.GET("/test", func(c *gin.Context) {

	})
	app.Run(":4000")
}
