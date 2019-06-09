package user

// User struct
type User struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name binding:"required""`
	Password string `json:"password binding:"required""`
}

type Repository interface {
	FindUserByID(id uint) (*User, error)
	FindUserByName(name string) (*User, error)
	SaveUser(user *User) error
	GetAllUsers() ([]User, error)
	DeleteUserByID(id uint) error
}
