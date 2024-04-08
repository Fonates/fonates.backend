package models

import (
	// "fonates.backend/pkg/utils"
	"errors"

	"fonates.backend/pkg/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Address    string `json:"address" gorm:"unique;not null"`
	Username   string `json:"username"`
	Status     string `json:"-"`
	AvatarUrl  string `json:"avatarUrl"`
}

func InitUser() User {
	return User{
		Status: "CREATED",
	}
}

func (u User) GetByAddress(store *gorm.DB, address string) (User, error) {
	result := store.Where("address = ?", address).First(&u)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return u, nil
	} else {
		return u, result.Error
	}
}

func (u User) Create(store *gorm.DB) (User, error) {
	u.Username = utils.GenerateUniqueID(12)
	return u, store.Create(&u).Error
}
