package service

import (
	"github.com/rs/zerolog/log"
	"github.com/zamblauskas/oxygen-control/config"
	"github.com/zamblauskas/oxygen-control/db"
	"github.com/zamblauskas/oxygen-control/models"
	"github.com/zamblauskas/oxygen-control/oxygen"
	"github.com/zamblauskas/oxygen-control/trigger"
)

func NewTriggerService(db *db.DB, oxygen *oxygen.Client, conf *config.Config) *TriggerService {
	return &TriggerService{db: db, oxygen: oxygen, conf: conf}
}

type TriggerService struct {
	db     *db.DB
	oxygen *oxygen.Client
	conf   *config.Config
}

func (s *TriggerService) LoadAndStartTriggers() error {
	triggers, err := s.db.GetAllTriggers()
	if err != nil {
		return err
	}

	for _, trigger := range triggers {
		switch trigger.GetType() {
		case models.TriggerTypeSchedule:
			s.startScheduleTrigger(trigger.(*models.ScheduleTrigger))
		case models.TriggerTypeFlic:
			s.startFlicTrigger(trigger.(*models.FlicTrigger))
		}
	}
	return nil
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
		err := s.db.AddFlicTrigger(v)
		if err != nil {
			return err
		}
		return s.startFlicTrigger(v)
	}
	return nil
}

func (s *TriggerService) GetTriggers() ([]models.Trigger, error) {
	return s.db.GetAllTriggers()
}

func (s *TriggerService) startScheduleTrigger(t *models.ScheduleTrigger) error {
	trigger.ScheduleStart(t.Hour, t.Minute, s.oxygen.Boost)
	return nil
}

func (s *TriggerService) startFlicTrigger(t *models.FlicTrigger) error {
	if s.conf.Flic != nil && s.conf.Flic.Enabled {
		callbacks := trigger.FlicCallbacks{
			OnButtonSingleClick: s.triggerActionToCallback(t.OnSingleClick),
			OnButtonDoubleClick: s.triggerActionToCallback(t.OnDoubleClick),
			OnButtonHold:        s.triggerActionToCallback(t.OnHold),
		}

		trigger.FlicListen(s.conf.Flic.ServerURL, t.Mac, callbacks)
	}
	return nil
}

func (s *TriggerService) triggerActionToCallback(action *models.TriggerAction) func() error {
	switch *action {
	case models.TriggerActionStartBoost:
		return s.oxygen.Boost
	case models.TriggerActionStopBoost:
		return s.oxygen.StopBoost
	}
	return nil
}
