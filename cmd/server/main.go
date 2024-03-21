package main

import (
	"os"

	"fonates.backend/pkg/api"
	"fonates.backend/pkg/databases/mariadb"
	"fonates.backend/pkg/routes"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	nodeEnv := os.Getenv("NODE_ENV")
	if "" == nodeEnv {
		nodeEnv = "development"
	}

	if err := godotenv.Load(".env." + nodeEnv); err != nil {
		panic(err)
	}

	config := api.Config{
		Mode:   nodeEnv,
		Host:   os.Getenv("SV_HOST"),
		Port:   os.Getenv("SV_PORT"),
		Router: routes.NewRouter("/api/v1"),
		Store: mariadb.MariaDB{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
			Database: os.Getenv("DB_NAME"),
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
