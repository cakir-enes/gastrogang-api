package recipe

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Recipe model
type Recipe struct {
	gorm.Model
	Name        string         `json:"name"`
	Steps       pq.StringArray `json:"steps" gorm:"type:varchar(100)[]"`
	Ingredients pq.StringArray `json:"ingredients" gorm:"type:varchar(100)[]"`
	Tags        pq.StringArray `json:"tags" gorm:"type:varchar(100)[]"`
	Details     string         `json:"details"`
	AuthorID    uint           `json:"authorid"`
}

type Repository interface {
	FindRecipeByID(id uint) (*Recipe, error)
	FindRecipesByAuthorID(id uint) ([]Recipe, error)
	SaveRecipe(recipe *Recipe) error
	DeleteRecipeByID(id uint) error
	UpdateRecipe(recipe *Recipe) error
	LikeRecipe(recId uint, userId uint) error
	DislikeRecipe(id uint, userId uint) error
	FindRecipeByTags(tags []string) ([]Recipe, error)
}
