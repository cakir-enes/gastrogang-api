package main

import (
	"gastrogang-api/pkg/httpd"
	"gastrogang-api/pkg/storage"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	repo := storage.NewInMemoryStorage()
	server := httpd.NewServer(repo, repo)
	server.Start()
}
