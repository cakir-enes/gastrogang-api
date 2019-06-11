package httpd

import (
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"
	"net/http"

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
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"msg": "pong"})
		})
	}
}
