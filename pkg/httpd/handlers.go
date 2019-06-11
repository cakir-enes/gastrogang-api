package httpd

import (
	"fmt"
	"gastrogang-api/pkg/storage"
	"gastrogang-api/pkg/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var failResp = func(msg string) interface{} {
	return gin.H{"status": "Fail", "message": msg}
}

var checkPassword = func(encrypted, normal string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(normal))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return false
	}
	return true
}

func registerUser(repo user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr user.User
		err := c.BindJSON(&usr)
		if err != nil {
			fmt.Println(err)
		}
		usr.HashPassword()
		err = repo.SaveUser(&usr)
		if err != nil {
			switch err {
			case storage.ConnectionFailed:
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			case storage.UserAlreadyExists:
				c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			}
			return
		}
		usr.GenerateToken()
		usr.Password = ""
		c.JSON(http.StatusOK, usr)
	}
}

func loginUser(repo user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usrJson user.User
		err := c.BindJSON(&usrJson)
		if err != nil {
			fmt.Println(err)
		}
		usr, err := repo.FindUserByName(usrJson.Name)
		if err != nil {
			switch err {
			case storage.UserDoesntExist:
				c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			case storage.ConnectionFailed:
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			default:
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			}
			return
		}
		match := checkPassword(usr.Password, usrJson.Password)
		if !match {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp("Invalid credentials"))
			return
		}
		c.Set("user", usr.ID)
		usr.GenerateToken()
		usr.Password = ""
		c.JSON(http.StatusOK, usr)
	}
}
