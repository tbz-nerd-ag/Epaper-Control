package types

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Epd struct {
	ID         string `json:"id"`
	Room       string `json:"room"`
	NightSleep bool   `json:"nightsleep"`
	Device     string `json:"device"`
}

var (
	EPDs   []Epd
	epdsMu sync.RWMutex
)

func (e *Epd) UnmarshalJSON(b []byte) error {
	type EpdAlias Epd
	aux := &struct {
		ID json.Number `json:"id"`
		*EpdAlias
	}{
		EpdAlias: (*EpdAlias)(e),
	}
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}
	e.ID = aux.ID.String() // 10 → "10"
	return nil
}

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
	go watchEPD(watcher)
}

func watchEPD(watcher *fsnotify.Watcher) {
	defer watcher.Close()
	var timer *time.Timer
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(100*time.Millisecond, func() {
					slog.Info("epd.json geändert, wird neu eingelesen")
					loadfromfileepd()
				})
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
	epdsMu.Lock()
	file, err := os.ReadFile("epd.json")
	if err != nil {
		slog.Error("Fehler beim Lesen der JSON", "error", err)
		os.Exit(1)
	}
	var config struct {
		EPD []Epd `json:"epd"`
	}
	if err := json.Unmarshal(file, &config); err != nil {
		slog.Error("Fehler beim Parsen der JSON", "error", err)
		os.Exit(1)
	}
	EPDs = config.EPD
	epdsMu.Unlock()
}

func GetRoomfromID(id string) string {
	epdsMu.RLock()
	defer epdsMu.RUnlock()

	for _, epd := range EPDs {
		if epd.ID == id {
			return epd.Room
		}
	}
	return ""
}

func SetNightSleep(id string, change bool) error {
	epdsMu.Lock()
	defer epdsMu.Unlock() // ← defer, nicht manuell vor WriteFile

	data, err := os.ReadFile("epd.json")
	if err != nil {
		return err
	}
	var config struct {
		EPD []Epd `json:"epd"`
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

	// Atomisches Schreiben: erst tmp, dann rename
	tmp := "epd.json.tmp"
	if err := os.WriteFile(tmp, formatted, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, "epd.json") // atomar auf Linux/macOS
}

func GetNightsleep(id string) (bool, error) {
	epdsMu.RLock()
	defer epdsMu.RUnlock()

	for _, entry := range EPDs {
		if entry.ID == id {
			return entry.NightSleep, nil
		}
	}
	return false, fmt.Errorf("id %s nicht gefunden", id)
}
