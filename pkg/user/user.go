package user

// User struct
type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Repository interface {
	FindUserByID(id uint) (*User, error)
	SaveUser(user *User) error
	GetAllUsers() ([]User, error)
	DeleteUserByID(id uint) error
}
