package database

import "GoAPIBackEnd/internal/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Users{})
	//DB.AutoMigrate(&models.Task{})
}
