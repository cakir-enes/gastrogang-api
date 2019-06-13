package httpd

import (
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"

	"github.com/gin-gonic/gin"
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
	s.router.Run(":8080")
}

func (s *server) initRoutes() {
	v1 := s.router.Group("/api/v1")
	v1.Use(user.JwtAuthentication())
	{
		v1.POST("/register", registerUser(s.userRepo))
		v1.POST("/login", loginUser(s.userRepo))

		v1.GET("/recipes", getAllRecipes(s.recipeRepo))
		v1.POST("/recipes", saveRecipe(s.recipeRepo))

		v1.GET("/recipes/:id", getRecipeByID(s.recipeRepo))
		v1.DELETE("/recipes/:id", deleteRecipeByID(s.recipeRepo))
		v1.PUT("/recipes/:id", updateRecipeByID(s.recipeRepo))
	}
}
