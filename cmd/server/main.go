package main

import (
	"os"

	"fonates.backend/pkg/api"
	"fonates.backend/pkg/databases/mariadb"
	"fonates.backend/pkg/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	config := api.Config{
		Host:   "127.0.0.1",
		Port:   "4035",
		Router: routes.NewRouter("/api/v1"),
		Store: mariadb.MariaDB{
			Host:     "localhost",
			Port:     "3306",
			Username: "root",
			Password: "24702470",
			Database: "fonates",
			// Username: "admin",
			// Password: "24702470",
		},
	}

	apiV1, errInit := api.NewApiServer(&config)
	if errInit != nil {
		panic(errInit)
	}

	if errStart := apiV1.Start(); errStart != nil {
		panic(errStart)
	}
}
