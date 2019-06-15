package postgres

import (
	"errors"
	"gastrogang-api/pkg/recipe"
)

func (s *Database) FindRecipeByID(id uint) (*recipe.Recipe, error) {
	var rec recipe.Recipe
	if s.db.First(&rec, id).RecordNotFound() {
		return nil, errors.New("RecipeDoesntExist")
	}
	return &rec, nil
}

func (s *Database) FindRecipesByAuthorID(id uint) ([]recipe.Recipe, error) {
	var recipes []recipe.Recipe
	s.db.Where("author_id = ?", id).Find(&recipes)
	return recipes, nil

}

func (s *Database) SaveRecipe(recipe *recipe.Recipe) error {
	s.db.Create(recipe)
	return nil
}

func (s *Database) DeleteRecipeByID(id uint) error {
	if s.db.Where("id = ?", id).Delete(&recipe.Recipe{}).RecordNotFound() {
		return errors.New("RecipeDoesntExist")
	}
	return nil
}

func (s *Database) UpdateRecipe(recipe *recipe.Recipe) error {
	s.db.Save(recipe)
	return nil
}
