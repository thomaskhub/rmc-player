package main

import (
	"fmt"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/gen2brain/go-mpv"

	"C"
)
import (
	"encoding/json"
)

type Player struct {
	playlist  Playlist
	decKeyMap map[string]string
	m         *mpv.Mpv
}

type ClientMessage struct {
	ArgsCount int
	Args      unsafe.Pointer
}

func NewPlayer(cfg Config) *Player {

	player := &Player{
		playlist: Playlist{},
		m:        mpv.New(),
	}

	player.m.SetPropertyString("input-default-bindings", "yes")
	player.m.SetPropertyString("input-vo-keyboard", "yes")
	player.m.SetPropertyString("input-conf", "./input.conf")

	if cfg.EnableLog {
		player.m.SetOptionString("log-file", "./rmc.log")
	}

	player.m.SetOptionString("osc", "no")
	player.m.SetOptionString("osd-in-seek", "msq-bar")
	player.m.SetOptionString("fs", "yes")
	player.m.SetOptionString("keep-open", "yes")
	player.m.SetOptionString("force-window", "yes")

	player.m.SetOptionString("script", "./lua/keybindings.lua")

	monitorPositions := GetMonitorPositions()

	if len(monitorPositions) > 1 {
		x := monitorPositions[1][0]
		y := monitorPositions[1][1]
		width := monitorPositions[1][2]
		height := monitorPositions[1][3]
		geometry := fmt.Sprintf("%dx%d+%d+%d", width, height, x, y)

		player.m.SetOptionString("geometry", geometry)
	}

	player.m.SetProperty("osd-duration", mpv.FormatInt64, 2000)
	player.m.ObserveProperty(0, "eof-reached", mpv.FormatFlag)

	player.decKeyMap = make(map[string]string)

	err := player.m.Initialize()
	if err != nil {
		panic(err)
	}

	//start the event loop
	go func() {
		for {
			e := player.m.WaitEvent(10000)
			switch e.EventID {
			case mpv.EventPropertyChange:
				prop := e.Property()

				switch prop.Name {
				case "eof-reached":
					if prop.Data != nil {
						if prop.Data.(int) == 1 {
							current := player.playlist.GetCurrent()
							if current != nil {
								if current.AutoPlayNext {
									player.PlayNext(false)
								}

							}
						}
					}
				}

			case mpv.EventClientMessage:
				tmp := (*ClientMessage)(e.Data)

				args := make([]string, tmp.ArgsCount)
				argsPtr := (*[1 << 30]*C.char)(unsafe.Pointer(tmp.Args))[:tmp.ArgsCount:tmp.ArgsCount]
				for i, argPtr := range argsPtr {
					args[i] = C.GoString(argPtr)
				}

				if args[0] == "keyevent" {
					if tmp.ArgsCount < 2 {
						//print warning that argument count does not match
						println("warn:: args count does not match")
						continue
					}

					switch args[1] {
					case "quit":
						os.Exit(0)
					case "enter":
						player.PlayNext(false)
					case "pageUp":
						player.PlayNext(false)
					case "pageDown":
						player.PlayNext(true)
					}
				}

			case mpv.EventShutdown:
				os.Exit(0)
			}
		}
	}()

	return player
}

func (p *Player) InitPlaylist(playlist []byte) {
	err := p.m.Command([]string{"loadfile", "./assets/screensaver.jpg"})
	if err != nil {
		panic(err)
	}

	p.playlist.jsonToPlaylist(playlist)
}

func (p *Player) PlaySingleFile(path string) {
	err := p.m.Command([]string{"loadfile", path})
	if err != nil {
		panic(err)
	}
}

func (p *Player) PlayListIndex(idx int) {
	p.playlist.Idx = idx
	p.PlayNext(false)
}

func (p *Player) PlayNext(direction bool) {
	var item *PlaylistItem
	if !direction {
		item = p.playlist.GetNext()
	} else {
		item = p.playlist.GetPrev()
	}

	if item == nil {
		return
	}

	p.m.SetOptionString("start", fmt.Sprintf("%d", item.Seek))
	extension := filepath.Ext(item.Path)

	//check if extension is a key in the p.decKeyMap
	if _, ok := p.decKeyMap[extension]; ok {
		data := fmt.Sprintf("protocol_whitelist=[crypto],decryption_key=%s", p.decKeyMap[extension])
		p.m.SetOptionString("demuxer-lavf-o", data)
	}

	//start playback of the pallist item path
	p.m.Command([]string{"loadfile", item.Path})
	p.m.SetOptionString("pause", "no")

}

func (p *Player) Close() {
	p.m.TerminateDestroy()
}

func (p *Player) SetDecMap(playlist []byte) {
	var data map[string]string
	json.Unmarshal(playlist, &data)

	for k, v := range data {
		p.decKeyMap[k] = v
	}
}
