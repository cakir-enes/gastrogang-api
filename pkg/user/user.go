package user

import (
	"github.com/jinzhu/gorm"
	"os"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
	Lifetime int    `json:"lifetime"`
}

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type Repository interface {
	FindUserByID(id uint) (*User, error)
	SaveUser(user *User) error
	GetAllUsers() ([]User, error)
	DeleteUserByID(id uint) error
	FindUserByName(name string) (*User, error)
}

func (u *User) HashPassword() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)
}

func (u *User) GenerateToken() {
	tk := &Token{UserId: u.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	u.Token = tokenString //Store the token in the response
}
