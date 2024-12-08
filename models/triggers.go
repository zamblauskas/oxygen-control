package models

import (
	"encoding/json"
	"fmt"
)

type TriggerType string

const (
	TriggerTypeSchedule TriggerType = "schedule"
	TriggerTypeFlic     TriggerType = "flic"
)

type TriggerAction string

const (
	TriggerActionStartBoost TriggerAction = "start_boost"
	TriggerActionStopBoost  TriggerAction = "stop_boost"
)

type Trigger interface {
	GetType() TriggerType
}

type FlicTrigger struct {
	ID            string         `json:"id"`
	Type          TriggerType    `json:"type"`
	Mac           string         `json:"mac" binding:"required"`
	OnSingleClick *TriggerAction `json:"on_single_click"`
	OnDoubleClick *TriggerAction `json:"on_double_click"`
	OnHold        *TriggerAction `json:"on_hold"`
}

func (t *FlicTrigger) GetType() TriggerType {
	return TriggerTypeFlic
}

type ScheduleTrigger struct {
	ID     string        `json:"id"`
	Type   TriggerType   `json:"type"`
	Hour   int           `json:"hour" binding:"required"`
	Minute int           `json:"minute" binding:"required"`
	Action TriggerAction `json:"action" binding:"required"`
}

func (t *ScheduleTrigger) GetType() TriggerType {
	return TriggerTypeSchedule
}

func ParseTriggerFromJson(data []byte) (Trigger, error) {
	var temp struct {
		Type TriggerType `json:"type"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}

	switch temp.Type {
	case TriggerTypeSchedule:
		var t ScheduleTrigger
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		return &t, nil
	case TriggerTypeFlic:
		var t FlicTrigger
		if err := json.Unmarshal(data, &t); err != nil {
			return nil, err
		}
		return &t, nil
	}

	return nil, fmt.Errorf("invalid trigger type")
}
