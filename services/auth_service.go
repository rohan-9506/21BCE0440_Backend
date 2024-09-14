package services

import (
	"file-sharing-system/models"
	"file-sharing-system/utils"
)

func Register(email, password string) error {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:        email,
		PasswordHash: hashedPassword,
	}

	result := models.DB.Create(&user)
	return result.Error
}

func Login(email, password string) (string, error) {
	var user models.User
	result := models.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return "", result.Error
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", result.Error // Change this to a specific error message if needed
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
