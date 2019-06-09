package httpd

import (
	"fmt"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func findUserRecipes(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		recipes, err := repo.FindRecipesByAuthor(username)
		if err != nil {
			c.JSON(http.StatusNotFound, fmt.Errorf("User with name: %s not found", username))
			return
		}
		c.JSON(http.StatusOK, recipes)
	}
}

func saveRecipe(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json recipe.Recipe
		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := repo.SaveRecipe(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}

func saveUser(repo user.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var json user.User
		if err := c.BindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if err := repo.SaveUser(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	}
}
