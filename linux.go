package main

import (
	"log"
	"runtime"
)

func InitLinux() {
	if runtime.GOOS == "linux" {
		log.Println("we are running on linux so init the linux system")

		if IsRaspberryPi() {
			//only for the pi we meddle with the sound. On laptop we keep it to the OS
			PulseSetVolume(20)
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
