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

func (s *Database) UpdateRecipe(newRecipe *recipe.Recipe) error {

	if s.db.Model(newRecipe).Updates(
		recipe.Recipe{
			Name: newRecipe.Name, Steps: newRecipe.Steps, Details: newRecipe.Details, Ingredients: newRecipe.Ingredients,
		}).RecordNotFound() {
		return errors.New("RecipeDoesntExist")
	}
	return nil
}

func (s *Database) LikeRecipe(recId uint, userId uint) error {
	var like recipe.Like
	if s.db.Where("recipe_id = ?", recId).Find(&like).RecordNotFound() {
		var rec recipe.Recipe
		if s.db.Where("id = ?", recId).Find(&rec).RecordNotFound() {
			return errors.New("RecipeDoesntExist")
		}
		like.Count = 1
		like.RecipeID = recId
		like.AuthorIDs = []int64{int64(userId)}
		s.db.Create(&like)
		return nil
	}
	for _, aId := range like.AuthorIDs {
		if aId == int64(userId) {
			return errors.New("User has already liked")
		}
	}
	like.Count += 1
	like.AuthorIDs = append(like.AuthorIDs, int64(userId))
	s.db.Save(&like)
	return nil
}

func (s *Database) DislikeRecipe(recId uint, userId uint) error {
	var like recipe.Like
	if s.db.Where("recipe_id = ?", recId).Find(&like).RecordNotFound() {
		var rec recipe.Recipe
		if s.db.Where("id = ?", recId).Find(&rec).RecordNotFound() {
			return errors.New("RecipeDoesntExist")
		}
		return errors.New("No one ever liked this")
	}
	for i, aId := range like.AuthorIDs {
		if aId == int64(userId) {
			like.Count -= 1
			like.AuthorIDs = remove(like.AuthorIDs, i)
			s.db.Save(&like)
			return nil
		}
	}
	return errors.New("User never liked this")
}

func remove(s []int64, i int) []int64 {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
