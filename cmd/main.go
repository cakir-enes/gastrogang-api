package main

import (
	"gastrogang-api/pkg/httpd"
	"gastrogang-api/pkg/storage"
)

func main() {
	repo := storage.NewInMemoryStorage()
	server := httpd.NewServer(repo, repo)
	server.Start()
}
