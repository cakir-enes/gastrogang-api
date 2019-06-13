package httpd

import (
	"errors"
	"fmt"
	"gastrogang-api/pkg/storage"
	"gastrogang-api/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var checkPassword = func(encrypted, normal string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(normal))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return false
	}
	return true
}

var extractIdFromCtx = func(c *gin.Context) (uint, error) {
	id, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": "FAIL", "msg": "Something went wrong"})
		return 0, errors.New("user doesnt exist in the context")
	}
	id, ok := id.(uint)
	if !ok {
		return 0, errors.New("Couldnt parse user id")
	}
	return id.(uint), nil
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
		var json user.User
		err := c.BindJSON(&json)
		if err != nil {
			fmt.Println(err)
		}
		usr, err := repo.FindUserByName(json.Name)
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
		match := checkPassword(usr.Password, json.Password)
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
