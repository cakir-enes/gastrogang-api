package httpd

import (
	"encoding/json"
	"gastrogang-api/pkg/recipe"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func appendLikeCount(recipes []recipe.Recipe, repo recipe.Repository) []map[string]interface{} {

	var resps []map[string]interface{}

	for _, rec := range recipes {
		like, err := repo.FindLikeOfRecipe(rec.ID)
		var resp map[string]interface{}
		recJson, _ := json.Marshal(rec)
		json.Unmarshal(recJson, &resp)
		if err != nil {
			resp["like"] = recipe.Like{Count: 0}
		} else {
			resp["like"] = like
		}
		resps = append(resps, resp)
	}
	return resps
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
			return
		}
		c.Status(http.StatusOK)
	}
}
