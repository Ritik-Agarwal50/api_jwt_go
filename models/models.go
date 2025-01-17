package models

import (
	"go-jwt/database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var GlobalDB *gorm.DB

type User struct {
	gorm.Model
	ID       int    `gorm:"pimaryKey"`
	Name     string `json:"name" binding:"require"`
	Email    string `json:"email" binding: "required" gorm: "unique"`
	Password string `json:"password" binding: "required"`
}

func (user *User) HashPassword(Password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CreateUserRecord() error {
	result := database.GlobalDB.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (user *User) CheckPassword(providePassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providePassword))
	if err != nil {
		return err
	}
	return nil
}
