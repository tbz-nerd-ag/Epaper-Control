package types

import (
	"encoding/json"
	"log"
	"os"
)

type config struct {
	Wartung            bool   `json:"wartung"`
	Wartung_sleep_time int    `json:"wartung_sleep_time"`
	Sleep_time         int    `json:"sleep_time"`
	Task_time_cron     string `json:"task_time_cron"`
}

var Config config

func Loadconfig() {
	file, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Fehler beim Lesen der JSON: ", err)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("Fehler beim Lesen der JSON: ", err)
	}
}
