package httpd

import (
	"fmt"
	"gastrogang-api/pkg/storage"
	"gastrogang-api/pkg/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

var failResp = func(msg string) interface{} {
	return gin.H{"status": "Fail", "message": msg}
}

func registerUser(repo user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var usr user.User
		err := c.BindJSON(&usr)
		if err != nil {
			fmt.Println(err)
		}
		usr.HashPwAndGenerateToken()
		err = repo.SaveUser(&usr)
		if err != nil {
			switch err {
			case storage.ConnectionFailed:
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			case storage.UsernameExists:
				c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			}
		}
		usr.Password = ""
		c.JSON(http.StatusOK, usr)
	}
}
