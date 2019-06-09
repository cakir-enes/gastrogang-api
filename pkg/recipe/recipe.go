package recipe

// Recipe model
type Recipe struct {
	ID       uint     `json:"id" binding:"required"`
	Name     string   `json:"name binding:"required""`
	Steps    []string `json:"steps binding:"required""`
	Details  string   `json:"details"`
	AuthorID uint     `json:"authorid" binding:"required"`
}

type Repository interface {
	FindRecipeByID(id uint) (*Recipe, error)
	FindRecipesByAuthor(name string) ([]Recipe, error)
	SaveRecipe(recipe *Recipe) error
	DeleteRecipeByID(id uint) error
	UpdateRecipe(recipe *Recipe) error
}
