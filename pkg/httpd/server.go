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
