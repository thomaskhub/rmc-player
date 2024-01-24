package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

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
