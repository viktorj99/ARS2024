package service

import (
	"fmt"
	"projekat/model"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

//todo: dodati metode za kreiranje, brisanje itd..

func (s ConfigService) Hello() {
	fmt.Println("Hello from config service")
}
