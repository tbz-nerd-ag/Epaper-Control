package types

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
)

type epd struct {
	ID         string `json:"id"`
	Room       string `json:"room"`
	NightSleep bool   `json:"nightsleep"`
}

var EPD epd

func Loadepd() {
	loadfromfileepd()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Fehler beim Erstellen des Watchers", "error", err)
		os.Exit(1)
	}
	if err := watcher.Add("epd.json"); err != nil {
		slog.Error("Fehler beim Hinzufügen der Datei", "error", err)
		os.Exit(1)
	}

	//subroutine that checks edits of epd.json
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
				slog.Info("epd.json geändert, wird neu eingelesen")
				loadFromFile()
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			slog.Error("Watcher-Fehler", "error", err)
		}
	}
}

func loadfromfileepd() {
	file, err := os.ReadFile("epd.json")
	if err != nil {
		slog.Error("Fehler beim Lesen der JSON", "error", err)
		os.Exit(1)
	}
	err = json.Unmarshal(file, &EPD)
	if err != nil {
		slog.Error("Fehler beim Lesen der JSON", "error", err)
		os.Exit(1)
	}
}

func GetRoomfromID(id string) string {
	idint, err := strconv.Atoi(id)
	if err != nil {
		return ""
	}

	data, err := os.ReadFile("epd.json")
	if err != nil {
		return ""
	}

	var config struct {
		EPD []struct {
			ID   int    `json:"id"`
			Room string `json:"room"`
		} `json:"epd"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return ""
	}

	for _, epd := range config.EPD {
		if epd.ID == idint {
			return epd.Room
		}
	}
	return ""
}

func SetNightSleep(id string, change bool) error {
	data, err := os.ReadFile("epd.json")
	if err != nil {
		return err
	}

	var config struct {
		EPD []epd `json:"epd"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	for i, entry := range config.EPD {
		if entry.ID == id {
			config.EPD[i].NightSleep = change
			break
		}
	}

	formatted, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("epd.json", formatted, 0644)
}

func GetNightsleep(id string) (bool, error) {
	data, err := os.ReadFile("epd.json")
	if err != nil {
		return false, err
	}

	var config struct {
		EPD []epd `json:"epd"`
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return false, err
	}

	for _, entry := range config.EPD {
		if entry.ID == id {
			return entry.NightSleep, nil
		}
	}

	return false, fmt.Errorf("id %s nicht gefunden", id)
}
