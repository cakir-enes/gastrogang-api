package httpd

import (
	"errors"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var failResp = func(msg string) interface{} {
	return gin.H{"status": "Fail", "message": msg}
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
		rec, err := findRecipeCheckAuthor(c, repo)
		if err != nil {
			return
		}
		c.JSON(http.StatusOK, rec)
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
			if err == storage.RecipeAlreadyExists {
				c.AbortWithStatusJSON(http.StatusUnprocessableEntity, failResp(err.Error()))
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		c.JSON(http.StatusOK, json)
	}
}

func deleteRecipeByID(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		recipe, err := findRecipeCheckAuthor(c, repo)
		if err != nil {
			return
		}
		repo.DeleteRecipeByID(recipe.ID)
		c.Status(http.StatusOK)
	}
}

func updateRecipeByID(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := findRecipeCheckAuthor(c, repo)
		if err != nil {
			return
		}
		var json recipe.Recipe
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			return
		}
		recipeId := uint(id)
		json.ID = recipeId
		err = c.BindJSON(&json)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			return
		}

		repo.UpdateRecipe(&json)
		c.Status(http.StatusOK)
	}
}

func likeRecipeByID(repo recipe.Repository) gin.HandlerFunc {
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
		err = repo.LikeRecipe(uint(id), userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
		}
		c.Status(http.StatusOK)
	}
}

func dislikeRecipeByID(repo recipe.Repository) gin.HandlerFunc {
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
		err = repo.DislikeRecipe(uint(id), userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
		}
		c.Status(http.StatusOK)
	}
}

func findRecipeCheckAuthor(c *gin.Context, repo recipe.Repository) (*recipe.Recipe, error) {
	userId, err := extractIdFromCtx(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
		return nil, err
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
		return nil, err
	}
	recipeId := uint(id)
	recipe, err := repo.FindRecipeByID(recipeId)
	if err != nil {
		if err == storage.ConnectionFailed {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			return nil, err
		}
		c.AbortWithStatusJSON(http.StatusNotFound, failResp("Recipe doesnt exist"))
		return nil, err
	}
	if recipe.AuthorID != userId {
		c.AbortWithStatusJSON(http.StatusForbidden, failResp("It's not even your recipe"))
		return nil, errors.New("Unauthorized")
	}
	return recipe, nil
}
