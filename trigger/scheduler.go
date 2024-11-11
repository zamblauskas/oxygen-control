package trigger

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Scheduler struct {
	trigger func() error
	hour    int
	minute  int
}

func NewScheduler(trigger func() error, hour, minute int) *Scheduler {
	return &Scheduler{
		trigger: trigger,
		hour:    hour,
		minute:  minute,
	}
}

func (s *Scheduler) Start() {
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day(), s.hour, s.minute, 0, 0, now.Location())
			if now.After(next) {
				next = next.Add(24 * time.Hour)
			}

			wait := next.Sub(now)
			log.Info().
				Int("hour", s.hour).
				Int("minute", s.minute).
				Msg("scheduled next trigger")

			time.Sleep(wait)

			if err := s.trigger(); err != nil {
				log.Error().Err(err).Msg("failed to execute scheduled trigger")
			}
		}
	}()
}
