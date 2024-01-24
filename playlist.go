package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type PlaylistItem struct {
	Path         string `json:"path"`
	AutoPlayNext bool   `json:"autoPlayNext"`
	PauseOnLast  bool   `json:"pauseOnLast"`
	Seek         int    `json:"seek"`
}

type Playlist struct {
	List []PlaylistItem
	Idx  int
}

func (p *Playlist) jsonToPlaylist(jsonString []byte) {
	p.List = make([]PlaylistItem, 0)
	tmpList := make([]PlaylistItem, 0)
	err := json.Unmarshal(jsonString, &tmpList)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	p.List = append(p.List, PlaylistItem{
		Path:         "./assets/screensaver.jpg",
		AutoPlayNext: false,
		PauseOnLast:  false,
		Seek:         0,
	})

	for _, item := range tmpList {
		//check if path exist if file does not exists print warning
		if _, err := os.Stat(item.Path); os.IsNotExist(err) {
			fmt.Printf("Warning: File does not exist: %s\n", item.Path)
			return
		}
	}

	for i, item := range tmpList {
		if i < len(tmpList)-1 && !item.AutoPlayNext && !item.PauseOnLast {
			item.AutoPlayNext = true
			p.List = append(p.List, item)
			p.List = append(p.List, PlaylistItem{
				Path:         "./assets/screensaver.jpg",
				AutoPlayNext: false,
				PauseOnLast:  false,
				Seek:         0,
			})
		} else {
			p.List = append(p.List, item)
		}
	}

	p.List = append(p.List, PlaylistItem{
		Path:         "./assets/screensaver.jpg",
		AutoPlayNext: false,
		PauseOnLast:  false,
		Seek:         0,
	})

	p.Idx = 0
}

func (p *Playlist) GetPlaylist() ([]PlaylistItem, int) {
	return p.List, p.Idx
}

func (p *Playlist) GetCurrent() *PlaylistItem {
	if p.Idx >= 0 && p.Idx < len(p.List) {
		return &p.List[p.Idx]
	} else {
		return nil
	}
}

func (p *Playlist) GetNext() *PlaylistItem {

	if p.Idx < len(p.List)-1 {
		p.Idx++
		return &p.List[p.Idx]
	} else {
		return nil
	}
}

func (p *Playlist) GetPrev() *PlaylistItem {

	if p.Idx-1 >= 0 {
		p.Idx--
		return &p.List[p.Idx]
	} else {
		return nil
	}
}
