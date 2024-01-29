package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/lawl/pulseaudio"
	"github.com/veandco/go-sdl2/sdl"
)

type Monitor struct {
	HasPrimary              bool
	HasSecondary            bool
	PrimaryResolutionList   []string
	SecondaryResolutionList []string
	PrimaryResolution       string
	SecondaryResolution     string
	PrimaryName             string
	SecondaryName           string
}

func GetAudioDevices() ([]string, error) {
	if runtime.GOOS == "linux" {
		// if runtime.GOOS == "linux" && (runtime.GOARCH == "arm" || runtime.GOARCH == "arm64") {
		client, err := pulseaudio.NewClient()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		defer client.Close()

		sinks, err := client.Sinks()
		if err != nil {
			log.Println(err)
			return nil, err
		}

		var ret []string
		for _, item := range sinks {
			ret = append(ret, item.Name)
		}

		keys := make(map[string]bool)
		list := []string{}
		for _, entry := range ret {
			if _, value := keys[entry]; !value {
				keys[entry] = true
				list = append(list, entry)
			}
		}

		return list, nil
	}
	return nil, fmt.Errorf("not running in linux")
}

func PlayerSetAudioDevice(device string) error {
	if runtime.GOOS == "linux" {
		client, err := pulseaudio.NewClient()
		if err != nil {
			log.Println(err)
			return err
		}
		defer client.Close()

		client.SetDefaultSink(device)
	}
	return nil
}
func GetMonitorCount() int {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	monitorCount, err := sdl.GetNumVideoDisplays()
	if err != nil {
		panic(err)
	}

	return monitorCount
}

func PulseSetVolume(vol int) error {
	if runtime.GOOS == "linux" {
		client, err := pulseaudio.NewClient()
		if err != nil {
			log.Println("Player::PoulseSetVolume::error --> ", err)
			return err
		}
		defer client.Close()

		client.SetVolume(float32(vol) / 100.0)
	}
	return nil
}

func PulseGetVolume() int {
	if runtime.GOOS == "linux" {
		client, err := pulseaudio.NewClient()
		if err != nil {
			log.Println("Player::PoulseSetVolume::error --> ", err)
			return -1
		}
		defer client.Close()

		pulseVol, err := client.Volume()
		if err != nil {
			return -1
		}

		var vol = int(pulseVol * 100)
		return vol
	}
	return 100
}

func GetMonitorPositions() [][]int32 {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	monitorCount, err := sdl.GetNumVideoDisplays()
	if err != nil {
		panic(err)
	}

	positions := make([][]int32, monitorCount)

	for i := 0; i < monitorCount; i++ {
		rect, err := sdl.GetDisplayBounds(i)
		if err != nil {
			panic(err)
		}
		positions[i] = []int32{rect.X, rect.Y, rect.W, rect.H}
	}

	return positions
}

func GetWindowsDrives() []string {
	var drives []string
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		if _, err := os.Stat(string(drive) + ":\\"); err == nil {
			drives = append(drives, string(drive)+":\\")
		}
	}
	return drives
}

// create a function that gets a system path string and return all directories path strings
// mp4, mp3 and wav files found in this path
func GetMediaFiles(path string) ([]FileListT, error) {
	var mediaFiles []FileListT

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		if strings.HasPrefix(file.Name(), ".") || PathIsInBlacklist(filepath.Join(path, file.Name())) {
			continue
		}

		if file.IsDir() {
			retPath := filepath.Join(path, file.Name())
			mediaFiles = append(mediaFiles, FileListT{IsDir: true, Path: retPath, Name: file.Name()})
		} else {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".mp4" || ext == ".mp3" || ext == ".wav" || ext == ".rmc" || strings.Contains(ext, ".sands") {
				retPath := filepath.Join(path, file.Name())
				mediaFiles = append(mediaFiles, FileListT{IsDir: false, Path: retPath, Name: file.Name()})
			}
		}
	}

	return mediaFiles, nil
}

func PathIsInBlacklist(path string) bool {
	switch runtime.GOOS {
	case "linux":
		//check if path starts with any item in the config.LinuxDirBlacklist
		for _, dir := range cfg.LinuxDirBlacklist {
			isInBlacklist := strings.HasPrefix(path, dir)
			if isInBlacklist {
				return true
			}
		}
	case "darwin":
		for _, dir := range cfg.DarwinDirBlacklist {
			isInBlacklist := strings.HasPrefix(path, dir)
			if isInBlacklist {
				return true
			}
		}
	case "windows":
		for _, dir := range cfg.WindowsDirBlacklist {
			isInBlacklist := strings.HasPrefix(path, dir)
			if isInBlacklist {
				return true
			}
		}
	}
	return false
}

var writeConfigMutex = &sync.Mutex{}

func UpdateConfigFile() {
	writeConfigMutex.Lock()
	defer writeConfigMutex.Unlock()

	configJSON, err := json.Marshal(cfg)
	if err == nil {
		os.WriteFile(configFile, configJSON, 0644)
	}
}

// func GetRandrSupportedResolutions() []string {
// 	out, err := exec.Command("xrandr").Output()
// 	if err != nil {
// 		return nil
// 	}

// 	var resolutions []string
// 	for _, line := range strings.Split(string(out), "\n") {
// 		parts := strings.Fields(line)

// 		if len(parts) == 2 {
// 			resolution := parts[0]
// 			resolutions = append(resolutions, resolution)
// 		}
// 	}

// 	//resolutions might contain duplicates remove them

// 	keys := make(map[string]bool)
// 	list := []string{}
// 	for _, entry := range resolutions {
// 		if _, value := keys[entry]; !value {
// 			keys[entry] = true
// 			tmp := strings.Split(entry, "x")

// 			width, err := strconv.Atoi(tmp[0])
// 			height, err1 := strconv.Atoi(tmp[1])
// 			if err == nil && err1 == nil && width >= 800 && height >= 600 {
// 				list = append(list, entry)
// 			}
// 		}
// 	}

// 	return list
// }

func GetRandrCurrentResolution() string {
	out, err := exec.Command("xrandr").Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "*") {
			parts := strings.Fields(line)
			return parts[0]
		}
	}

	return "1920x1080"
}

const (
	MONITOR_PRIMARY   = 0
	MONITOR_SECONDARY = 1
)

func GetRandrResolutionsFromLines(data []string) ([]string, string) {
	currentRes := ""
	var resolutions []string
	for _, line := range data {
		parts := strings.Fields(line)

		if len(parts) >= 2 && strings.Contains(parts[0], "x") {
			if strings.Contains(parts[1], "*") {
				currentRes = parts[0]
			}
			resolution := parts[0]

			//split resolution by x and only append if width >= 800 and height >= 600
			tmp := strings.Split(resolution, "x")
			width, err := strconv.Atoi(tmp[0])
			height, err1 := strconv.Atoi(tmp[1])
			if err == nil && err1 == nil && width >= 800 && height >= 600 {
				resolutions = append(resolutions, resolution)
			}

			// order the list by highest to smallest width values
			sort.Slice(resolutions, func(i, j int) bool {
				w1, _ := strconv.Atoi(strings.Split(resolutions[i], "x")[0])
				w2, _ := strconv.Atoi(strings.Split(resolutions[j], "x")[0])
				return w1 > w2
			})

			//now ensure to remove duplicate entries
			keys := make(map[string]bool)
			list := []string{}
			for _, entry := range resolutions {
				if _, value := keys[entry]; !value {
					keys[entry] = true
					list = append(list, entry)
				}
			}
			resolutions = list
		}
	}

	return resolutions, currentRes
}

func GetRandrMonitorName(line string) string {
	parts := strings.Fields(line)
	return parts[0]
}

func GetRandrMonitorDetails() *Monitor {
	out, err := exec.Command("xrandr").Output()
	if err != nil {
		return nil
	}
	monitorLines := make([]int, 0)
	data := make([]string, 0)

	for i, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, "connected") {
			monitorLines = append(monitorLines, i)
		}
		data = append(data, line)
	}

	monitor := Monitor{HasPrimary: false, HasSecondary: false}

	//so now we now can access the monitors
	if len(monitorLines) == 1 {
		//we only have primary no secondary
		monitor.HasPrimary = true
		monitor.PrimaryResolutionList, monitor.PrimaryResolution = GetRandrResolutionsFromLines(data)
		monitor.HasSecondary = false
		monitor.SecondaryResolution = ""
		monitor.SecondaryResolutionList = []string{""}
		monitor.PrimaryName = GetRandrMonitorName(data[monitorLines[0]])
	} else if len(monitorLines) > 1 {
		//we have primary and secondary
		monitor.HasPrimary = true
		monitor.HasSecondary = true

		prim := data[0:monitorLines[1]]
		sec := data[monitorLines[0]+1:]
		monitor.PrimaryResolutionList, monitor.PrimaryResolution = GetRandrResolutionsFromLines(prim)
		monitor.SecondaryResolutionList, monitor.SecondaryResolution = GetRandrResolutionsFromLines(sec)

		monitor.PrimaryName = GetRandrMonitorName(data[monitorLines[0]])
		monitor.SecondaryName = GetRandrMonitorName(data[monitorLines[1]])
	} else {
		return nil
	}

	return &monitor

}

func SetRandrMonitorResolution(primary bool, resolution string) error {
	monitor := GetRandrMonitorDetails()
	if monitor == nil {
		return errors.New("no monitor details available")
	}

	name := ""
	if primary {
		name = monitor.PrimaryName
	} else {
		name = monitor.SecondaryName
	}

	// Use the xrandr command to set the resolution
	cmd := exec.Command("xrandr", "--output", name, "--mode", resolution)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set resolution for %s monitor: %v", name, err)
	}

	return nil
}
