package httpd

func (s *server) initRoutes() {
	v1 := s.router.Group("/api/v1")
	{
		v1.GET("/user/:username/recipes", findUserRecipes(s.recipeRepo))
		v1.POST("/user/:username/recipes", saveRecipe(s.recipeRepo))
		v1.POST("/users", saveUser(s.userRepo))
	}
}
