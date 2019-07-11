package main

import (
	"gastrogang-api/pkg/httpd"
	"gastrogang-api/pkg/storage/postgres"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	repo := postgres.NewPgDB()

	defer repo.Close()

	server := httpd.NewServer(repo, repo)
	server.Start()
}
