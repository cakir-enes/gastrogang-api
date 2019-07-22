package httpd

import (
	"gastrogang-api/pkg/recipe"
	"github.com/gin-gonic/gin"
	"net/http"
)

func uploadPhotos(repo recipe.Repository) gin.HandlerFunc {
	type Req struct {
		Photos []recipe.Photo
	}
	return func(c *gin.Context) {
		recipe, err := findRecipeCheckAuthor(c, repo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		var req Req
		c.BindJSON(&req)
		for _, p := range req.Photos {
			p.RecipeID = recipe.ID
			err := repo.SavePhoto(&p)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			}
		}
		c.Status(200)
	}
}

func getPhotosByID(repo recipe.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		recipe, err := findRecipeCheckAuthor(c, repo)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, failResp(err.Error()))
			return
		}
		photos, err := repo.GetPhotosByID(recipe.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, failResp(err.Error()))
			return
		}
		c.JSON(http.StatusOK, photos)
	}
}
