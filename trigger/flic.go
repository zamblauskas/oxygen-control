package trigger

import (
	"github.com/rs/zerolog/log"

	"github.com/zamblauskas/flic-go-ws-client/flic"
)

type FlicCallbacks struct {
	OnButtonSingleClick func() error
	OnButtonDoubleClick func() error
	OnButtonHold        func() error
}

func FlicListen(flicServerURL string, buttonBluetoothAddress string, callbacks FlicCallbacks) {
	go func() {
		client, err := flic.NewClient(flicServerURL)
		if err != nil {
			log.Error().Err(err).Msg("failed to create Flic client")
			return
		}

		defer client.Close()

		client.OnButton = func(event flic.ButtonEvent) {
			switch event.ClickType {
			case flic.ButtonSingleClick:
				log.Info().Msg("Flic button single click")
				if callbacks.OnButtonSingleClick != nil {
					if err := callbacks.OnButtonSingleClick(); err != nil {
						log.Error().Err(err).Msg("failed to execute on button single click")
					}
				}
			case flic.ButtonDoubleClick:
				log.Info().Msg("Flic button double click")
				if callbacks.OnButtonDoubleClick != nil {
					if err := callbacks.OnButtonDoubleClick(); err != nil {
						log.Error().Err(err).Msg("failed to execute on button double click")
					}
				}
			case flic.ButtonHold:
				log.Info().Msg("Flic button hold")
				if callbacks.OnButtonHold != nil {
					if err := callbacks.OnButtonHold(); err != nil {
						log.Error().Err(err).Msg("failed to execute on button hold")
					}
				}
			}
		}

		if err := client.Connect(buttonBluetoothAddress); err != nil {
			log.Error().Err(err).Msg("failed to connect to Flic button")
			return
		}

		if err := client.Listen(); err != nil {
			log.Error().Err(err).Msg("failed to listen for Flic button events")
			return
		}

		log.Info().Msg("listening for Flic button events")
	}()
}
