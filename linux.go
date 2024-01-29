package main

import (
	"log"
	"runtime"
)

func InitLinux() {
	if runtime.GOOS == "linux" {
		log.Println("we are running on linux so init the linux system")
		//Audio Setup
		PulseSetVolume(100)
		devices, err := GetAudioDevices()
		if err == nil { //set first output if nothing is stored
			var hasDefault = false
			for _, dev := range devices {
				if dev == cfg.AudioOutput {
					PlayerSetAudioDevice(cfg.AudioOutput)
					hasDefault = true
					break
				}
			}

			if !hasDefault {
				cfg.AudioOutput = devices[0]
				PlayerSetAudioDevice(cfg.AudioOutput)
			}
		}

		//Resolution Setup
		mon := GetRandrMonitorDetails()
		if mon != nil {
			if cfg.PrimaryResolution == "" {
				cfg.PrimaryResolution = mon.PrimaryResolution
				SetRandrMonitorResolution(true, cfg.PrimaryResolution)
			} else {
				SetRandrMonitorResolution(true, cfg.PrimaryResolution)
			}

			if cfg.SecondaryResolution == "" {
				cfg.SecondaryResolution = mon.SecondaryResolution
				SetRandrMonitorResolution(false, cfg.SecondaryResolution)
			} else {
				SetRandrMonitorResolution(false, cfg.SecondaryResolution)
			}
		}

	}
}
