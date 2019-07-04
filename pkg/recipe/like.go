package recipe

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Like struct {
	gorm.Model
	Count     uint `json:"count"`
	RecipeID  uint
	AuthorIDs pq.Int64Array `gorm:"type:integer[]"`
}
