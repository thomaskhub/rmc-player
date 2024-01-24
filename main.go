package main

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port                string   `json:"port"`
	EnableLog           bool     `json:"enableLog"`
	LinuxDirBlacklist   []string `json:"linuxDirBlacklist"`
	DarwinDirBlacklist  []string `json:"darwinDirBlacklist"`
	WindowsDirBlacklist []string `json:"windowsDirBlacklist"`
}

var cfg Config

func main() {

	//please load the config.json file and unmarshal into Config struct
	//check if config file exist
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		cfg.DarwinDirBlacklist = []string{"/boot", "/lost+found", "/opt", "/root", "/etc", "/sys", "/usr", "/bin", "/sbin", "/dev", "/proc", "/run", "/snap", "/tmp", "/srv", "/var", "/cdrom"}
		cfg.EnableLog = true
		cfg.LinuxDirBlacklist = []string{"/boot", "/lost+found", "/opt", "/root", "/etc", "/sys", "/usr", "/bin", "/sbin", "/dev", "/proc", "/run", "/snap", "/tmp", "/srv", "/var", "/cdrom"}
		cfg.WindowsDirBlacklist = []string{}
	} else {
		file, err := os.ReadFile("config.json")
		if err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		if err := json.Unmarshal(file, &cfg); err != nil {
			log.Fatalf("Error unmarshaling config: %v", err)
		}
	}

	//in case no port is set used a default setting
	if cfg.Port == "" {
		cfg.Port = "8880"
	}

	player := NewPlayer(cfg)
	player.InitPlaylist([]byte("[]"))

	server := Server{
		Player: player,
	}

	server.Run(cfg.Port)
}
