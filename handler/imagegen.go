package handler

import (
	"Control/untis"
	"encoding/json"
	"os"
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

func PrepareJSON(raum string) {
	data, err := os.ReadFile("untis/cache/" + raum + ".json")
	if err != nil {
		untis.Get_data(raum)
	}

	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return
	}

	today := time.Now().Format("2006-01-02")

	filtered := []Lesson{}
	for _, lesson := range resp.Lessons {
		if lesson.Date == today {
			filtered = append(filtered, lesson)
		}
	}
	resp.Lessons = filtered
}
