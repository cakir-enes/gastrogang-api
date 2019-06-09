package recipe

// Recipe model
type Recipe struct {
	ID       uint     `json:"id"`
	Name     string   `json:"name"`
	Steps    []string `json:"steps"`
	Details  string   `json:"details"`
	AuthorID uint     `json:"authorid"`
}

type Repository interface {
	FindRecipeByID(id uint) (*Recipe, error)
	FindRecipesByAuthor(name string) ([]Recipe, error)
	SaveRecipe(recipe *Recipe) error
	DeleteRecipeByID(id uint) error
	UpdateRecipe(recipe *Recipe) error
}
