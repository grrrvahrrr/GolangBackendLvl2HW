package apigin

import (
	"GoBeLvl2/redis"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiGin struct {
	*gin.Engine
	Redis *redis.RedisStore
}

func NewApiGin(redis *redis.RedisStore) *ApiGin {
	apigin := &ApiGin{
		Redis: redis,
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.POST("/register", apigin.register)
	router.POST("/auth", apigin.auth)

	apigin.Engine = router

	return apigin
}

func (ag *ApiGin) register(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	authCode, err := ag.Redis.CacheRegister(c, login, password)
	if err != nil {
		c.Error(err)
		return
	}
	if authCode == "user_already_registred" {
		c.String(http.StatusBadRequest, "User Login: %s has already registered, please send your authenticatioon code to /auth", login)
		return
	}

	fmt.Printf("User Login: %s, AuthCode: %s", login, authCode)
	c.String(http.StatusOK, "User Login: %s, AuthCode: %s", login, authCode)

}

func (ag *ApiGin) auth(c *gin.Context) {
	authCode := c.PostForm("authCode")

	err := ag.Redis.CacheAuth(c, authCode)
	if err != nil {
		c.Error(err)
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Authentication successful!")
	c.String(http.StatusOK, "Authentication successful!")
}
