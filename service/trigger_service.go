package service

import (
	"github.com/rs/zerolog/log"
	"github.com/zamblauskas/oxygen-control/db"
	"github.com/zamblauskas/oxygen-control/models"
	"github.com/zamblauskas/oxygen-control/oxygen"
	"github.com/zamblauskas/oxygen-control/trigger"
)

func NewTriggerService(db *db.DB) *TriggerService {
	return &TriggerService{db: db}
}

type TriggerService struct {
	db     *db.DB
	oxygen *oxygen.Client
}

func (s *TriggerService) AddTrigger(t models.Trigger) error {
	switch v := t.(type) {
	case *models.ScheduleTrigger:
		log.Info().Msgf("Adding schedule trigger: %v", v)
		err := s.db.AddScheduleTrigger(v)
		if err != nil {
			return err
		}
		return s.startScheduleTrigger(v)
	case *models.FlicTrigger:
		log.Info().Msgf("Adding flic trigger: %v", v)
		return s.db.AddFlicTrigger(v)
	}
	return nil
}

func (s *TriggerService) startScheduleTrigger(t *models.ScheduleTrigger) error {
	trigger.ScheduleStart(t.Hour, t.Minute, s.oxygen.Boost)
	return nil
}
