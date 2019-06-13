package httpd

import (
	"errors"
	"fmt"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/storage"
	"gastrogang-api/pkg/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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

func getAllRecipes(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := extractIdFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		recipes, err := repo.FindRecipesByAuthorID(id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		c.JSON(http.StatusOK, recipes)
	}
}

func getRecipeByID(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := extractIdFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			return
		}
		recipeId := uint(id)
		recipe, err := repo.FindRecipeByID(recipeId)
		if err != nil {
			if err == storage.ConnectionFailed {
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			return
		}
		if recipe.AuthorID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, failResp("Recipe is not public"))
			return
		}
		c.JSON(http.StatusOK, recipe)
	}
}

func saveRecipe(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := extractIdFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		var json recipe.Recipe
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			return
		}
		json.AuthorID = id
		err = repo.SaveRecipe(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		c.JSON(http.StatusOK, json)
	}
}

func deleteRecipeByID(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := extractIdFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			return
		}
		recipeId := uint(id)
		recipe, err := repo.FindRecipeByID(recipeId)
		if err != nil {
			if err == storage.ConnectionFailed {
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp("Recipe doesnt exist"))
			return
		}
		if recipe.AuthorID != userId {
			c.AbortWithStatusJSON(http.StatusUnauthorized, failResp("It's not even your recipe"))
			return
		}
		repo.DeleteRecipeByID(recipeId)
		c.Status(http.StatusOK)
	}
}
