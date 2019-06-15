package postgres

import (
	"gastrogang-api/pkg/storage"
	"gastrogang-api/pkg/user"
)

func (s *Database) FindUserByID(id uint) (*user.User, error) {
	var u user.User
	if s.db.First(&u, id).RecordNotFound() {
		return nil, storage.UserDoesntExist
	}
	return &u, nil
}

func (s *Database) SaveUser(newUser *user.User) error {
	var u user.User
	if !s.db.Where("name = ?", newUser.Name).First(&u).RecordNotFound() {
		return storage.UserAlreadyExists
	}
	err := s.db.Create(newUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Database) GetAllUsers() ([]user.User, error) {
	var users []user.User
	err := s.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *Database) DeleteUserByID(id uint) error {
	return s.db.Where("id = ?", id).Delete(&user.User{}).Error
}

func (s *Database) FindUserByName(name string) (*user.User, error) {
	var u user.User
	if s.db.Where("name = ?", name).First(&u).RecordNotFound() {
		return nil, storage.UserDoesntExist
	}

	return &u, nil
}
