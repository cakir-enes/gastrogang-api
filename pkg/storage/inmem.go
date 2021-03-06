package storage

import (
	"fmt"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"
)

type inMemory struct {
	recipes   []recipe.Recipe
	users     []user.User
	userCount uint
	recCount  uint
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

func (s *inMemory) FindRecipesByAuthorID(id uint) ([]recipe.Recipe, error) {
	recipes := []recipe.Recipe{}
	for _, r := range s.recipes {
		if r.AuthorID == id {
			recipes = append(recipes, r)
		}
	}
	return recipes, nil
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
	recipes, err := s.FindRecipesByAuthorID(recipe.AuthorID)
	if err != nil {
		return err
	}
	for _, rec := range recipes {
		if rec.Name == recipe.Name {
			return RecipeAlreadyExists
		}
	}
	recipe.ID = s.recCount + 1
	s.recCount = s.recCount + 1
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
	var idx = -1
	for i, r := range s.recipes {
		if r.ID == recipe.ID {
			idx = i
		}
	}
	s.recipes[idx].Name = recipe.Name
	s.recipes[idx].Steps = recipe.Steps
	s.recipes[idx].Details = recipe.Details
	s.recipes[idx].Ingredients = recipe.Ingredients

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
	return nil, UserDoesntExist
}
func (s *inMemory) SaveUser(user *user.User) error {
	_, err := s.FindUserByName(user.Name)
	if err != UserDoesntExist {
		return UserAlreadyExists
	}
	user.ID = s.userCount + 1
	s.userCount = s.userCount + 1
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
