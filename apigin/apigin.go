package apigin

import (
	"GoBeLvl2/elastic"
	"GoBeLvl2/enteties"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiGin struct {
	*gin.Engine
	storage elastic.PostStorage
}

func NewApiGin(storage elastic.PostStorage) *ApiGin {
	apigin := &ApiGin{
		storage: storage,
	}

	router := gin.New()

	router.Use(gin.Logger())
	router.POST("/register", apigin.register)
	router.POST("/finduser", apigin.findUser)
	router.POST("/searchuser", apigin.searchUser)

	apigin.Engine = router

	return apigin
}

func (ag *ApiGin) register(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")
	age := c.PostForm("age")
	sex := c.PostForm("sex")

	user := enteties.User{
		Login:    login,
		Password: password,
		Age:      age,
		Sex:      sex,
	}

	err := ag.storage.Insert(c, user)
	if err != nil {
		c.Error(err)
	}

	fmt.Printf("User Login: %s", login)
	c.String(http.StatusOK, "User Login: %s", login)

}

func (ag *ApiGin) findUser(c *gin.Context) {

	login := c.PostForm("login")

	user, err := ag.storage.FindOne(c, login)
	if err != nil {
		c.Error(err)
	}

	fmt.Printf("User found: %v", user)
	c.String(http.StatusOK, "User found: %s", user)

}

func (ag *ApiGin) searchUser(c *gin.Context) {

	key := c.PostForm("key")
	value := c.PostForm("value")

	users, err := ag.storage.SearchUser(c, key, value)
	if err != nil {
		c.Error(err)
	}

	fmt.Printf("User found: %v", users)
	c.String(http.StatusOK, "User found: %v", users)

}
