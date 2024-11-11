package main

import (
	"github.com/rs/zerolog/log"
	"github.com/zamblauskas/oxygen-control/config"
	"github.com/zamblauskas/oxygen-control/logger"
	"github.com/zamblauskas/oxygen-control/oxygen"
	"github.com/zamblauskas/oxygen-control/trigger"
)

func main() {
	logger.Setup()

	conf, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	client := oxygen.NewClient(conf.OxygenURL)

	for _, schedule := range conf.Schedules {
		scheduler := trigger.NewScheduler(client.Boost, schedule.Hour, schedule.Minute)
		scheduler.Start()
	}

	err = client.Boost()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to trigger oxygen boost")
	}

	log.Info().Msg("Oxygen control started")
	log.Info().Msg("Press Ctrl+C to stop")
	select {}
}
