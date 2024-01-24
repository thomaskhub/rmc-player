package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/gen2brain/go-mpv"
)

type Server struct {
	Player            *Player
	EventTimer        *time.Ticker
	EventTimerStarted bool
}

func (s *Server) Run(port string) {
	http.HandleFunc("/secret", s.secretHandler)
	http.HandleFunc("/playlist", s.playlistHandler)
	http.HandleFunc("/play", s.play)
	http.HandleFunc("/alive", s.aliveHandler)
	http.HandleFunc("/files", s.filesHandler)
	http.HandleFunc("/cmd", s.commandHandler)
	http.HandleFunc("/mediaInfo", s.mediaInfoHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

type FileListT struct {
	IsDir bool   `json:"isDir"`
	Path  string `json:"path"`
	Name  string `json:"name"`
}

func (s *Server) commandHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// Read the command from the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Execute the command
		result := make(map[string]interface{})
		json.Unmarshal(body, &result)

		switch result["cmd"] {
		case "selectPlaylist":
			params := result["params"].(map[string]interface{})
			idx := int(params["idx"].(float64))
			s.Player.PlayListIndex(idx - 1)
		case "play":
			s.Player.m.SetOptionString("pause", "no")
		case "playPause":
			s.Player.m.Command([]string{"cycle", "pause"})
		case "pause":
			s.Player.m.SetOptionString("pause", "yes")
		case "mute":
			s.Player.m.Command([]string{"cycle", "mute"})
		case "next":
			s.Player.PlayNext(false)
		case "previous":
			s.Player.PlayNext(true)
		case "seek":
			params := result["params"].(map[string]interface{})
			time := fmt.Sprintf("%d", int(params["time"].(float64)))
			println(time)
			s.Player.m.Command([]string{"seek", time, "absolute"})
		}
		w.WriteHeader(http.StatusOK)

	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) mediaInfoHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		timePos, _ := s.Player.m.GetProperty("time-pos", mpv.FormatInt64)
		duration, _ := s.Player.m.GetProperty("duration", mpv.FormatInt64)
		volume := s.Player.m.GetPropertyString("volume")
		muted := s.Player.m.GetPropertyString("mute")
		paused := s.Player.m.GetPropertyString("pause")

		playlist, idx := s.Player.playlist.GetPlaylist()

		if playlist == nil {
			playlist = []PlaylistItem{}
		}

		if duration == nil {
			duration = 0
		}

		var timePosVal int
		switch timePos.(type) {
		case int:
			timePosVal = timePos.(int)
		case int64:
			timePosVal = int(timePos.(int64))
		default:
			timePosVal = 0
		}

		data := map[string]interface{}{
			"timePos":  timePosVal,
			"playlist": playlist,
			"idx":      idx,
			"duration": duration,
			"volume":   volume,
			"muted":    muted == "yes",
			"paused":   paused == "yes",
		}

		jsonData, _ := json.Marshal(data)
		w.Write(jsonData)
	}
}

func (s *Server) filesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	//get the directory specified in the url, if not path is specified use
	//the root directory of the os

	directory := r.URL.Query().Get("path")
	var root []FileListT
	if directory == "" {
		switch runtime.GOOS {
		case "windows":
			//for windows directory will be a list of all the drives found in the system
			wDrives := GetWindowsDrives()
			for _, drive := range wDrives {
				root = append(root, FileListT{IsDir: true, Path: drive, Name: drive})
			}
			json.NewEncoder(w).Encode(root)
			return
		}

		//search for all direc
		mFiles, _ := GetMediaFiles("/")
		parentDir := strings.Split(directory, string(os.PathSeparator))
		parentDir = parentDir[:len(parentDir)-1]
		directory = strings.Join(parentDir, string(os.PathSeparator))
		mFiles = append([]FileListT{{IsDir: true, Path: directory, Name: "..."}}, mFiles...)
		json.NewEncoder(w).Encode(mFiles)
		return
	}

	mFiles, _ := GetMediaFiles(directory)
	parentDir := strings.Split(directory, string(os.PathSeparator))
	parentDir = parentDir[:len(parentDir)-1]
	directory = strings.Join(parentDir, string(os.PathSeparator))

	mFiles = append([]FileListT{{IsDir: true, Path: directory, Name: "..."}}, mFiles...)
	json.NewEncoder(w).Encode(mFiles)
}

func (s *Server) secretHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	// Get the data from the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	s.Player.SetDecMap(body)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) playlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	// Get the data from the body of the request
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	s.Player.InitPlaylist(body)
}

func (s *Server) play(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body.", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var tmp PlaylistItem
	json.Unmarshal(body, &tmp)

	//check if tmp.path ends with .rmc
	if strings.HasSuffix(tmp.Path, ".rmc") {
		//read the file from tm.path and decode the json
		file, err := os.ReadFile(tmp.Path)
		if err != nil {
			http.Error(w, "Failed to read file.", http.StatusInternalServerError)
			return
		}

		var playlist []PlaylistItem
		err = json.Unmarshal(file, &playlist)
		if err != nil {
			http.Error(w, "Failed to decode JSON.", http.StatusInternalServerError)
			return
		}

		for _, item := range playlist {
			path := filepath.Base(item.Path)
			fullPath := filepath.Join(tmp.Path, path)
			item.Path = fullPath

			if _, err := os.Stat(item.Path); os.IsNotExist(err) {
				http.Error(w, "Failed to find file.", http.StatusInternalServerError)
				return
			}
		}

		data, _ := json.Marshal(playlist)
		s.Player.InitPlaylist(data)

	} else {
		//Directly play the file
		data, _ := json.Marshal([]PlaylistItem{tmp})
		s.Player.InitPlaylist(data)
		s.Player.PlayNext(false)
	}
}

func (s *Server) aliveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusMethodNotAllowed)
		return
	}
}
