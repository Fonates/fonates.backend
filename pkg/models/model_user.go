package models

import (
	// "fonates.backend/pkg/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Address    string `json:"address" gorm:"unique;not null"`
	Username   string `json:"username"`
	Status     string `json:"status"`
	AvatarUrl  string `json:"avatarUrl"`
}

func InitUser() User {
	return User{
		Status: "CREATED",
	}
}

func (u User) Create(store *gorm.DB) (User, error) {
	// u.Username = utils.GenerateUniqueID(12)
	return u, store.Create(&u).Error
}
