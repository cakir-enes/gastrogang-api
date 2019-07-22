package httpd

import (
	"context"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
)

type server struct {
	router     *gin.Engine
	userRepo   user.Repository
	recipeRepo recipe.Repository
}

func NewServer(userRepo user.Repository, recipeRepo recipe.Repository) *server {
	server := server{router: gin.Default(), userRepo: userRepo, recipeRepo: recipeRepo}
	server.initRoutes()
	return &server
}

func (s *server) Start() {
	err := s.router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Errorf(context.Background(), "Cant bind port\n")
	}
}

func (s *server) initRoutes() {
	s.router.Static("/swagger", "cmd/swaggerui")
	s.router.Use(corsCfg())
	s.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger")
		c.Abort()
	})
	v1 := s.router.Group("/api/v1")
	v1.Use(corsCfg())
	v1.Use(user.JwtAuthentication())
	{
		v1.POST("/register", registerUser(s.userRepo))
		v1.POST("/login", loginUser(s.userRepo))

		v1.GET("/recipes", getAllRecipes(s.recipeRepo))
		v1.POST("/recipes", saveRecipe(s.recipeRepo))
		v1.GET("/search", getRecipeByTags(s.recipeRepo))

		v1.GET("/recipes/:id", getRecipeByID(s.recipeRepo))
		v1.DELETE("/recipes/:id", deleteRecipeByID(s.recipeRepo))
		v1.PUT("/recipes/:id", updateRecipeByID(s.recipeRepo))

		v1.POST("/recipes/:id/like", likeRecipeByID(s.recipeRepo))
		v1.POST("/recipes/:id/dislike", dislikeRecipeByID(s.recipeRepo))
		v1.POST("/recipes/:id/photo", uploadPhotos(s.recipeRepo))
		v1.GET("/recipes/:id/photo", getPhotosByID(s.recipeRepo))
		v1.POST("/recipes/:id/toggle-publicity", togglePublicity(s.recipeRepo))
	}
}

func corsCfg() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Content-Length"},
		MaxAge:           12 * time.Hour,
		AllowCredentials: false,
	})
}
