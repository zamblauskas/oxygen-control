package trigger

import (
	"time"

	"github.com/rs/zerolog/log"
)

func ScheduleStart(hour, minute int, callback func() error) {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			wait := next.Sub(now)
			log.Info().
				Int("hour", hour).
				Int("minute", minute).
				Msg("scheduled next trigger")

			time.Sleep(wait)

			if err := callback(); err != nil {
				log.Error().Err(err).Msg("failed to execute scheduled trigger")
			}
		}
	}()
}
