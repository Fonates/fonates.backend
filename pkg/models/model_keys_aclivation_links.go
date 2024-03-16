package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KeysActivationLink struct {
	Id             uint      `json:"id" gorm:"primaryKey, autoIncrement, not null"`
	Key            uuid.UUID `json:"key" gorm:"type:uuid"`
	Status         string    `json:"status" gorm:"not null"`
	DonationLinkID uint      `json:"donation_link_id" gorm:"not null"`
	DonationLink   DonationLink
}

func InitKeysActivation(linkId uint) KeysActivationLink {
	return KeysActivationLink{
		Status:         "INACTIVE",
		Key:            uuid.New(),
		DonationLinkID: linkId,
	}
}

func (k KeysActivationLink) Create(store *gorm.DB) error {
	return store.Create(&k).Error
}
