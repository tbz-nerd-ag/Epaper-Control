package types

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/fsnotify/fsnotify"
)

type config struct {
	Wartung            bool     `json:"wartung"`
	Wartung_sleep_time int      `json:"wartung_sleep_time"`
	Task_time_cron     []string `json:"task_time_cron"`
	InfluxToken        string   `json:"influxtoken"`
	Log_Filename       string   `json:"log_filename"`
	Log_Max_Age        int      `json:"log_max_age"`
	Log_Max_Backups    int      `json:"log_max_backups"`
	Log_compress       bool     `json:"log_compress"`
}

var Config config

func Loadconfig() {
	loadFromFile()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		slog.Error("Fehler beim Erstellen des Watchers", "error", err)
		os.Exit(0)
	}
	if err := watcher.Add("config.json"); err != nil {
		slog.Error("Fehler beim Hinzufügen der Datei", "error", err)
		os.Exit(1)
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
				slog.Info("config.json geändert, wird neu eingelesen")
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

func loadFromFile() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		slog.Error("Fehler beim Lesen der JSON", "error", err)
		os.Exit(1)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		slog.Error("Fehler beim Lesen der JSON", "error", err)
		os.Exit(1)
	}
}

func SaveConfig() error {
	data, err := json.MarshalIndent(Config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("config.json", data, 0644)
}
