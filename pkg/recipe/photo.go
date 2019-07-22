package recipe

import (
	"github.com/jinzhu/gorm"
)

type Photo struct {
	gorm.Model
	RecipeID uint   `json:"recipeid"`
	Type     string `json:"type"`
	Img      string `json:"img"`
}
