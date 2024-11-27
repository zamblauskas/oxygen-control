package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/zamblauskas/oxygen-control/api"
	"github.com/zamblauskas/oxygen-control/config"
	"github.com/zamblauskas/oxygen-control/db"
	"github.com/zamblauskas/oxygen-control/logger"
	"github.com/zamblauskas/oxygen-control/oxygen"
	"github.com/zamblauskas/oxygen-control/service"
	"github.com/zamblauskas/oxygen-control/trigger"
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

	triggerService := service.NewTriggerService(db)
	triggerHandler := api.NewTriggerHandler(triggerService)

	router := gin.Default()
	router.POST("/triggers", triggerHandler.AddTrigger)
	router.Run("127.0.0.1:8512")

	client := oxygen.NewClient(conf.OxygenURL)

	for _, schedule := range conf.Schedules {
		trigger.ScheduleStart(schedule.Hour, schedule.Minute, client.Boost)
	}

	if conf.Flic != nil && conf.Flic.Enabled {
		callbacks := trigger.FlicCallbacks{
			OnButtonSingleClick: client.Boost,
			OnButtonDoubleClick: client.StopBoost,
		}

		trigger.FlicListen(conf.Flic.ServerURL, conf.Flic.ButtonBluetoothAddr, callbacks)
	}

	log.Info().Msg("Oxygen control started")
	log.Info().Msg("Press Ctrl+C to stop")
	select {}
}
