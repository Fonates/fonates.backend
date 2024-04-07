package mariadb

import (
	"fonates.backend/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connector interface {
	Connect() (*gorm.DB, error)
	Migration(db *gorm.DB) error
}

type MariaDB struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func InitMariaDB() Connector {
	return &MariaDB{
		Host:     "",
		Port:     "",
		Username: "",
		Password: "",
		Database: "",
	}
}

func (m *MariaDB) Connect() (*gorm.DB, error) {
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (m *MariaDB) Migration(db *gorm.DB) error {
	return db.AutoMigrate(
		models.User{},
		models.DonationLink{},
		models.KeysActivationLink{},
	)
}
