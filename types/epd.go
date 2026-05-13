package types

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type epd struct {
	ID   string `json:"id"`
	Room string `json:"room"`
}

var EPD epd

func Loadepd() {
	loadfromfileepd()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Fehler beim Erstellen des Watchers: ", err)
	}
	if err := watcher.Add("epd.json"); err != nil {
		log.Fatal("Fehler beim Hinzufügen der Datei: ", err)
	}

	//subroutine that checks edits of config.json
	go watchEPD(watcher)
}

func watchEPD(watcher *fsnotify.Watcher) {
	defer watcher.Close()
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				log.Println("config.json wurde geändert, wird neue eingelesen ...")
				loadFromFile()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher-Fehler:", err)
		}
	}
}

func loadfromfileepd() {
	file, err := os.ReadFile("epd.json")
	if err != nil {
		log.Fatal("Fehler beim Lesen der JSON: ", err)
	}
	err = json.Unmarshal(file, &EPD)
	if err != nil {
		log.Fatal("Fehler beim Lesen der JSON: ", err)
	}
}
