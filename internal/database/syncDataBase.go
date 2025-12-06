package database

import "GoAPIBackEnd/internal/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.UsersLifttyn{})
	//DB.AutoMigrate(&models.Task{})
}
