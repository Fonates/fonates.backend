package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	LINK_KEY_INACTIVE = "INACTIVE"
	LINK_KEY_ACTIVE   = "ACTIVE"
)

type KeysActivationLink struct {
	ID             uint         `json:"id" gorm:"primaryKey, autoIncrement, not null"`
	Key            uuid.UUID    `json:"key" gorm:"type:uuid"`
	Status         string       `json:"status" gorm:"not null"`
	DonationLinkID uint         `json:"donationLinkId" gorm:"not null"`
	DonationLink   DonationLink `json:"-"`
}

func InitKeysActivation(linkId uint) KeysActivationLink {
	return KeysActivationLink{
		Status:         LINK_KEY_INACTIVE,
		Key:            uuid.New(),
		DonationLinkID: linkId,
	}
}

func (k KeysActivationLink) Create(store *gorm.DB) error {
	return store.Create(&k).Error
}

func (k KeysActivationLink) Activate(store *gorm.DB) error {
	return store.Model(&k).Update("status", "ACTIVE").Error
}

func (k KeysActivationLink) GetByLinkID(store *gorm.DB) (KeysActivationLink, error) {
	return k, store.Where("donation_link_id = ?", k.DonationLinkID).First(&k).Error
}
