package database

import (
	"GoAPIBackEnd/internal/models"
	_ "github.com/gin-gonic/gin"
)

func GetUser(name string) (*models.Users, error) {
	var user models.Users

	err := DB.Where("name = ?", name).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func PostUser(user models.Users) error {

	err := DB.Create(user)
	if err != nil {
		return err.Error
	}

	return nil
}
