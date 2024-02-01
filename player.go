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
	mpv       *mpv.Mpv
}

type ClientMessage struct {
	ArgsCount int
	Args      unsafe.Pointer
}

func NewPlayer(cfg Config) *Player {

	player := &Player{
		playlist: Playlist{},
		mpv:      mpv.New(),
	}

	player.mpv.SetPropertyString("input-default-bindings", "yes")
	player.mpv.SetPropertyString("input-vo-keyboard", "yes")
	player.mpv.SetPropertyString("input-conf", "./input.conf")

	if cfg.EnableLog {
		player.mpv.SetOptionString("log-file", "./rmc.log") //TODO: this should come from the config file
	}

	player.mpv.SetOptionString("osc", "no")
	player.mpv.SetOptionString("osd-in-seek", "msq-bar")
	player.mpv.SetOptionString("fs", "yes")
	player.mpv.SetOptionString("keep-open", "yes")
	player.mpv.SetOptionString("force-window", "yes")

	player.mpv.SetOptionString("script", "./lua/keybindings.lua")

	fullscreenGeometry := GetFullscreenGeometry()
	player.mpv.SetOptionString("geometry", fullscreenGeometry)

	player.mpv.SetProperty("osd-duration", mpv.FormatInt64, 2000)
	player.mpv.ObserveProperty(0, "eof-reached", mpv.FormatFlag)

	player.decKeyMap = make(map[string]string)

	err := player.mpv.Initialize()
	if err != nil {
		panic(err)
	}

	//start the event loop
	go func() {
		for {
			e := player.mpv.WaitEvent(10000)
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
	err := p.mpv.Command([]string{"loadfile", ScreensaverImage})
	// err := p.m.Command([]string{"loadfile", "./assets/screensaver.jpg"})
	if err != nil {
		panic(err)
	}

	p.playlist.jsonToPlaylist(playlist)
}

func (p *Player) PlaySingleFile(path string) {
	err := p.mpv.Command([]string{"loadfile", path})
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

	p.mpv.SetOptionString("start", fmt.Sprintf("%d", item.Seek))
	extension := filepath.Ext(item.Path)

	//check if extension is a key in the p.decKeyMap
	if _, ok := p.decKeyMap[extension]; ok {
		data := fmt.Sprintf("protocol_whitelist=[crypto],decryption_key=%s", p.decKeyMap[extension])
		p.mpv.SetOptionString("demuxer-lavf-o", data)
	}

	//start playback of the pallist item path
	p.mpv.Command([]string{"loadfile", item.Path})
	p.mpv.SetOptionString("pause", "no")

}

func (p *Player) Close() {
	p.mpv.TerminateDestroy()
}

func (p *Player) SetDecMap(playlist []byte) {
	var data map[string]string
	json.Unmarshal(playlist, &data)

	for k, v := range data {
		p.decKeyMap[k] = v
	}
}
