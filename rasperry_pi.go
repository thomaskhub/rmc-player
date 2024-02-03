package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func InitRaspberryPi() {
	fmt.Printf("runtime.GOOS: %v\n", runtime.GOOS)
	if IsRaspberryPi() {
		log.Println("we are running on PI so init the pi system system")

		//ensure the volume output is full power. The actual volume will be
		//controlled with the mixer
		PulseSetVolume(100)

		//disable the screensaver using xset -d :0 s off command
		cmd := exec.Command("xset", "-d", ":0", "s", "off")
		err := cmd.Run()
		if err != nil {
			log.Println(err)
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
