package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Config struct {
	Port                string   `json:"port"`
	EnableLog           bool     `json:"enableLog"`
	LinuxDirBlacklist   []string `json:"linuxDirBlacklist"`
	DarwinDirBlacklist  []string `json:"darwinDirBlacklist"`
	WindowsDirBlacklist []string `json:"windowsDirBlacklist"`
	AssetPath           string   `json:"assetPath"`
	AudioOutput         string   `json:"audioOutput"`
	PrimaryResolution   string   `json:"PrimaryResolution"`
	SecondaryResolution string   `json:"SecondaryResolution"`
}

var cfg Config

var ScreensaverImage string
var configFile string

func main() {
	log.Println("load the configuration")
	//read flag -c to specify a different config file path
	flag.StringVar(&configFile, "c", "config.json", "config file path")
	flag.Parse()

	if _, err := os.Stat(configFile); !os.IsNotExist(err) {
		file, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		if err := json.Unmarshal(file, &cfg); err != nil {
			log.Fatalf("Error unmarshaling config: %v", err)
		}
	}

	//if any of the parameters in the config file is not set use some defaults
	if cfg.DarwinDirBlacklist == nil || len(cfg.DarwinDirBlacklist) == 0 {
		cfg.DarwinDirBlacklist = []string{"/boot", "/lost+found", "/opt", "/root", "/etc", "/sys", "/usr", "/bin", "/sbin", "/dev", "/proc", "/run", "/snap", "/tmp", "/srv", "/var", "/cdrom"}
	}

	//if any of the parameters in the config file is not set use some defaults
	if cfg.LinuxDirBlacklist == nil || len(cfg.LinuxDirBlacklist) == 0 {
		cfg.LinuxDirBlacklist = []string{"/boot", "/lost+found", "/opt", "/root", "/etc", "/sys", "/usr", "/bin", "/sbin", "/dev", "/proc", "/run", "/snap", "/tmp", "/srv", "/var", "/cdrom"}
	}

	//if any of the parameters in the config file is not set use some defaults
	if cfg.WindowsDirBlacklist == nil || len(cfg.WindowsDirBlacklist) == 0 {
		cfg.WindowsDirBlacklist = []string{}
	}

	//in case no port is set used a default setting
	if cfg.Port == "" {
		cfg.Port = "8880"
	}

	//if assetPath is not define set it to ./assets
	if cfg.AssetPath == "" {
		cfg.AssetPath = "./assets"
	}

	ScreensaverImage = cfg.AssetPath + "/screensaver.jpg"

	// Prepare the audio system if we are running linux
	// needed for Raspberry PI. For Linux PC its nice to have feature
	InitLinux()

	//
	// Start the player
	//
	log.Println("create player")
	player := NewPlayer(cfg)
	player.InitPlaylist([]byte("[]"))

	server := Server{
		Player: player,
	}

	log.Println("starting webserver")
	server.Run(cfg.Port)
}
