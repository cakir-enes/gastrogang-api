package recipe

// Recipe model
type Recipe struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name" binding:"required"`
	Steps    []string `json:"steps" binding:"required"`
	Details  string   `json:"details" binding:"required"`
	AuthorID uint     `json:"authorid"`
}

type Repository interface {
	FindRecipeByID(id uint) (*Recipe, error)
	FindRecipesByAuthorID(id uint) ([]Recipe, error)
	SaveRecipe(recipe *Recipe) error
	DeleteRecipeByID(id uint) error
	UpdateRecipe(recipe *Recipe) error
}
