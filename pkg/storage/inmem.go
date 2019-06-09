package storage

import (
	"fmt"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"
)

type inMemory struct {
	recipes []recipe.Recipe
	users   []user.User
}

func NewInMemoryStorage() *inMemory {
	storage := inMemory{recipes: []recipe.Recipe{}, users: []user.User{}}
	return &storage
}

func (s *inMemory) FindRecipeByID(id uint) (*recipe.Recipe, error) {
	for _, r := range s.recipes {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("Recipe with id: %d not found", id)
}

func (s *inMemory) FindRecipesByAuthor(name string) ([]recipe.Recipe, error) {
	usr, err := s.FindUserByName(name)
	if err != nil {
		return nil, err
	}
	recipes := []recipe.Recipe{}
	for _, recipe := range s.recipes {
		if recipe.AuthorID == usr.ID {
			recipes = append(recipes, recipe)
		}
	}
	return recipes, nil
}

func (s *inMemory) SaveRecipe(recipe *recipe.Recipe) error {
	s.recipes = append(s.recipes, *recipe)
	return nil
}

func (s *inMemory) DeleteRecipeByID(id uint) error {
	newRecipes := []recipe.Recipe{}
	for _, r := range s.recipes {
		if r.ID != id {
			newRecipes = append(newRecipes, r)
		}
	}
	s.recipes = newRecipes
	return nil
}

func (s *inMemory) UpdateRecipe(recipe *recipe.Recipe) error {
	rec, err := s.FindRecipeByID(recipe.ID)
	if err != nil {
		return err
	}
	rec.Name = recipe.Name
	rec.Steps = recipe.Steps
	rec.Details = recipe.Details
	rec.AuthorID = recipe.AuthorID
	return nil
}

func (s *inMemory) FindUserByID(id uint) (*user.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("User with id: %d not found", id)
}

func (s *inMemory) FindUserByName(name string) (*user.User, error) {
	for _, user := range s.users {
		if user.Name == name {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("User with name: %s not found", name)
}
func (s *inMemory) SaveUser(user *user.User) error {
	s.users = append(s.users, *user)
	return nil
}
func (s *inMemory) GetAllUsers() ([]user.User, error) {
	return s.users, nil
}

func (s *inMemory) DeleteUserByID(id uint) error {
	newUsers := []user.User{}
	for _, user := range s.users {
		if user.ID != id {
			newUsers = append(newUsers, user)
		}
	}
	s.users = newUsers
	return nil
}
