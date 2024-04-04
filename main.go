package main

import (
	"projekat/repository"
	"projekat/service"
)

func main() {
	repo := repository.NewConfigConsulRepository()
	service := service.NewConfigService(repo)
	service.Hello()
}
