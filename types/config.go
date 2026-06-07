package types

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
)

type config struct {
	Wartung            bool     `json:"wartung"`
	Wartung_sleep_time int      `json:"wartung_sleep_time"`
	Sleep_time         int      `json:"sleep_time"`
	Task_time_cron     []string `json:"task_time_cron"`
	InfluxToken        string   `json:"influxtoken"`
}

var Config config

func Loadconfig() {
	loadFromFile()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Fehler beim Erstellen des Watchers: ", err)
	}
	if err := watcher.Add("config.json"); err != nil {
		log.Fatal("Fehler beim Hinzufügen der Datei: ", err)
	}

	//subroutine that checks edits of config.json
	go watchConfig(watcher)
}

func watchConfig(watcher *fsnotify.Watcher) {
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

func loadFromFile() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Fehler beim Lesen der JSON: ", err)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("Fehler beim Lesen der JSON: ", err)
	}
}

func SaveConfig() error {
	data, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}
