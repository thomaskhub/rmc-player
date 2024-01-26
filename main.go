package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
}

var cfg Config

var ScreensaverImage string

func main() {
	//read flag -c to specify a different config file path
	var configFile string
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

	fmt.Printf("cfg.AssetPath: %v\n", cfg.AssetPath)
	fmt.Printf("configFile: %v\n", configFile)

	ScreensaverImage = cfg.AssetPath + "/screensaver.jpg"

	player := NewPlayer(cfg)
	player.InitPlaylist([]byte("[]"))

	server := Server{
		Player: player,
	}

	server.Run(cfg.Port)
}
