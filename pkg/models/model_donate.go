package models

import (
	"gorm.io/gorm"
)

type Donate struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Hash           string `json:"hash" gorm:"unique;not null"`
	Amount         uint64 `json:"amount" gorm:"not null"`
	Comment        string `json:"comment,omitempty" gorm:"default:null"`
	Username       string `json:"username" gorm:"not null"`
	DonationLinkID uint   `json:"linkId" gorm:"not null"`
	UserID         uint   `json:"userId" gorm:"not null"`
}

func InitDonate() Donate {
	return Donate{}
}

func (d *Donate) Create(store *gorm.DB) error {
	return store.Create(d).Error
}
