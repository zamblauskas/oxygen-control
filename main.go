package main

import (
	"github.com/rs/zerolog/log"
	"github.com/zamblauskas/oxygen-control/api"
	"github.com/zamblauskas/oxygen-control/config"
	"github.com/zamblauskas/oxygen-control/db"
	"github.com/zamblauskas/oxygen-control/logger"
	"github.com/zamblauskas/oxygen-control/oxygen"
	"github.com/zamblauskas/oxygen-control/service"
)

func main() {
	logger.Setup()

	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db, err := db.NewDB(conf.DBPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database")
	}

	client := oxygen.NewClient(conf.OxygenURL)

	triggerService := service.NewTriggerService(db, client, conf)

	triggerHandler := api.NewTriggerHandler(triggerService)
	server := api.NewServer(triggerHandler)
	server.SetupRoutes()
	server.Start("127.0.0.1:8512")

	log.Info().Msg("Oxygen control started")
	log.Info().Msg("Press Ctrl+C to stop")
	select {}
}
