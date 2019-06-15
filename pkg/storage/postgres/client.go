package postgres

import (
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/user"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
)

type Database struct {
	db *gorm.DB
}

func NewPgDB() *Database {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatal("Cant connect to DB ", err.Error())
	}
	db = db.AutoMigrate(&user.User{}, &recipe.Recipe{})
	return &Database{db}
}

func (s *Database) Close() {
	s.db.Close()
}
