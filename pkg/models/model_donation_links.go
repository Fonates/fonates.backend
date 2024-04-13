package models

import (
	"errors"
	"regexp"

	"fonates.backend/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	LINK_INACTIVE = "INACTIVE"
	LINK_ACTIVE   = "ACTIVE"
	LINK_BLOCKED  = "BLOCKED"
)

type DonationLink struct {
	ID     uint      `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	UUID   uuid.UUID `json:"uuid" gorm:"unique;not null"`
	Name   string    `json:"custom_name"`
	Link   string    `json:"link" gorm:"unique;not null"`
	Status string    `json:"status"`
	UserID uint      `json:"-"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`
}

func InitDonationLink() DonationLink {
	return DonationLink{
		Status: LINK_INACTIVE,
	}
}

func (d DonationLink) Validate() bool {
	regexpLink := regexp.MustCompile(`^https:\/\/fonates\.com\/donates\/.*`)
	// regexpUsername := regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
	// isValidAddress := utils.ValidateTonAddress(d.Address)
	// && regexpUsername.MatchString(d.Username) && isValidAddress
	return regexpLink.MatchString(d.Link)
}

func (d DonationLink) Create(store *gorm.DB) (DonationLink, error) {
	d.Name = "dt" + utils.GenerateUniqueID(16)
	return d, store.Create(&d).Error
}

func (d DonationLink) GetById(store *gorm.DB, id string) (*DonationLink, error) {
	result := store.Where("id = ?", id).First(&d)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &d, nil
}

func (d DonationLink) Activate(store *gorm.DB) error {
	return store.Model(&d).Update("status", "ACTIVE").Error
}
