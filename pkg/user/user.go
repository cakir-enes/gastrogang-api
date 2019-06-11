package user

import "github.com/dgrijalva/jwt-go"

// User struct
type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Service interface {
	FindUserByID(id uint) (*User, error)
	SaveUser(user *User) error
	GetAllUsers() ([]User, error)
	DeleteUserByID(id uint) error
}
