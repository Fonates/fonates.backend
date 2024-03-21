package api

import (
	"net/http"

	"fonates.backend/pkg/databases/mariadb"
	"fonates.backend/pkg/handlers"
	"fonates.backend/pkg/routes"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type API interface {
	Start() error
}

type Config struct {
	Mode   string
	Host   string
	Port   string
	Router routes.Router
	Store  mariadb.MariaDB
}

type api struct {
	Store    *gorm.DB
	Router   *routes.Router
	Server   *http.Server
	Handlers *handlers.Handlers
}

func NewApiServer(config *Config) (api, error) {
	store, err := config.Store.Connect()
	if err != nil {
		log.Error().Err(err).Msg("Connect to MariaDB")
		return api{}, err
	}

	log.Info().Msg("Connected to MariaDB")

	if errMigrate := config.Store.Migration(store); errMigrate != nil {
		log.Error().Err(errMigrate).Msg("Migration MariaDB")
		return api{}, err
	}

	log.Info().Msg("Migrated MariaDB")

	handlers := handlers.NewHandlers(store, config.Mode)

	return api{
		Store:  store,
		Router: &config.Router,
		Server: &http.Server{
			Handler: config.Router.InitRoutes(*handlers),
			Addr:    config.Host + ":" + config.Port,
		},
	}, nil
}

func (s api) Start() error {
	log.Info().Msg("Starting API server")
	if err := s.Server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
