package main

import (
	"gastrogang-api/pkg/httpd"
	"gastrogang-api/pkg/recipe"
	"gastrogang-api/pkg/storage"
	"gastrogang-api/pkg/user"
)

func main() {
	str := storage.NewInMemoryStorage()
	str.SaveUser(&user.User{ID: 1, Name: "TestUser", Password: "1234"})
	str.SaveRecipe(&recipe.Recipe{ID: 1, Name: "TestRecipe", AuthorID: 1})
	server := httpd.NewServer(str, str)
	server.Start()

}
