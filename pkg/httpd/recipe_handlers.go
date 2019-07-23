package httpd

import (
	"errors"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
		c.JSON(http.StatusOK, appendLikeCount(recipes, repo))
	}
}

func getRecipeByID(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
			return
		}
		recipeId := uint(id)
		rec, err := repo.FindRecipeByID(recipeId)
		if err != nil {
			if err == storage.ConnectionFailed {
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp("Something went wrong"))
				return
			}
			c.AbortWithStatusJSON(http.StatusNotFound, failResp("Recipe Not Found"))
			return
		}
		// If recipe is public other logged in users can view it.
		if !rec.IsPublic {
			userId, err := extractIdFromCtx(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
				return
			}
			if rec.AuthorID != userId {
				c.AbortWithStatusJSON(http.StatusForbidden, failResp("Bad User."))
				return
			}
		}
		c.JSON(http.StatusOK, appendLikeCount([]recipe.Recipe{*rec}, repo)[0])
	}
}

func getRecipeByTags(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tags := c.Request.URL.Query()["tag"]
		//fmt.Printf("Tags %v len: %d \n", tags, len(tags))
		recipes, err := repo.FindRecipeByTags(tags)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		userId, err := extractIdFromCtx(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		userRecs := []recipe.Recipe{}
		for _, rec := range recipes {
			if rec.AuthorID == userId {
				userRecs = append(userRecs, rec)
			}
		}
		c.JSON(http.StatusOK, appendLikeCount(userRecs, repo))
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

func togglePublicity(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		rec, err := findRecipeCheckAuthor(c, repo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		isPublic, err := repo.TogglePublicity(rec.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			return
		}
		c.JSON(http.StatusOK, gin.H{"isPublic": isPublic})
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
