package models

import (
	"errors"
	"regexp"

	"fonates.backend/pkg/utils"
	"gorm.io/gorm"
)

type DonationLink struct {
	gorm.Model `json:"-"`
	Link       string `json:"link" gorm:"unique;not null"`
	Username   string `json:"username" gorm:"unique;not null"`
	Address    string `json:"address" gorm:"unique;not null"`
	Status     string `json:"status"`
}

func InitDonationLink() DonationLink {
	return DonationLink{
		Status: "INACTIVE",
	}
}

func (d DonationLink) Validate() bool {
	regexpLink := regexp.MustCompile(`^https:\/\/fonates\.com\/donates\/.*`)
	regexpUsername := regexp.MustCompile(`^[a-zA-Z0-9_]{3,16}$`)
	isValidAddress := utils.ValidateTonAddress(d.Address)
	return regexpLink.MatchString(d.Link) && regexpUsername.MatchString(d.Username) && isValidAddress
}

func (d DonationLink) Create(store *gorm.DB) (DonationLink, error) {
	return d, store.Create(&d).Error
}

func (d DonationLink) GetByAddress(store *gorm.DB, address string) (*DonationLink, error) {
	result := store.Where("address = ?", address).First(&d)
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
