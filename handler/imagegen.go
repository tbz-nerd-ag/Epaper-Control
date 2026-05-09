package handler

import (
	"Control/untis"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Lesson struct {
	Anzahl      int    `json:"anzahl"`
	Classroom   string `json:"classroom"`
	Code        string `json:"code"`
	Date        string `json:"date"`
	EndTime     string `json:"end_time"`
	Klasse      string `json:"klasse"`
	StartTime   string `json:"start_time"`
	Subject     string `json:"subject"`
	Teacher     string `json:"teacher"`
	RoomChanged bool   `json:"room_changed,omitempty"`
}

type Response struct {
	Lessons []Lesson `json:"lessons"`
	Room    string   `json:"room"`
}

type Room struct {
	Room string `json:"room"`
}

const outputDir = "handler/cache/"

func PrepareJSON(raum string) {
	data, err := os.ReadFile("untis/room.json")
	if err != nil {
		panic(err)
	}

	var rooms []Room
	if err := json.Unmarshal(data, &rooms); err != nil {
		panic(err)
	}

	// Ordner erstellen falls nicht vorhanden
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Println("Ordner erstellen Err:", err)
		return
	}

	for _, r := range rooms {
		data, err := os.ReadFile("untis/cache/" + r.Room + ".json")
		if err != nil {
			untis.Get_data(raum)
		}

		var resp Response
		if err := json.Unmarshal(data, &resp); err != nil {
			return
		}

		today := now().Format("2006-01-02")

		filtered := []Lesson{}
		for _, lesson := range resp.Lessons {
			if lesson.Date == today {
				filtered = append(filtered, lesson)
			}
		}
		resp.Lessons = filtered

		// JSON formatieren
		formatted, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			fmt.Printf("Marshal Fehler %s: %v\n", r.Room, err)
			continue
		}

		// Pfad zusammensetzen: "data/roomName.json"
		filename := filepath.Join(outputDir, r.Room+".json")
		err = os.WriteFile(filename, formatted, 0644)
		if err != nil {
			fmt.Println("Datei schreiben Err:", err)
			return
		}
	}

}

func now() time.Time {
	return time.Date(2026, 5, 11, 0, 0, 0, 0, time.Local)
}
