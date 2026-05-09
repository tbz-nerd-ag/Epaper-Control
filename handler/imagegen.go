package handler

import (
	"Control/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	RoomChanged bool   `json:"room_changed"`
}

type Response struct {
	Lessons []Lesson `json:"lessons"`
	Room    string   `json:"room"`
}

type ImageResponse struct {
	Status string `json:"status"`
	Image  string `json:"image"` // "0xff, 0xff, 0xff, ..."
}

const outputDir = "handler/cache/"
const hexDir = "handler/image_hex/"

func PrepareJSON() {
	data, err := os.ReadFile("untis/room.json")
	if err != nil {
		panic(err)
	}

	var rooms []types.Room
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
			continue
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

func Getpicturehex() {
	data, err := os.ReadFile("untis/room.json")
	if err != nil {
		panic(err)
	}

	var rooms []types.Room
	if err := json.Unmarshal(data, &rooms); err != nil {
		panic(err)
	}

	//create folder for image hex json
	err = os.MkdirAll(hexDir, 0755)
	if err != nil {
		fmt.Println("Ordner erstellen Err:", err)
		return
	}

	for _, r := range rooms {
		formatted, err := os.ReadFile(filepath.Join(outputDir, r.Room+".json"))
		if err != nil {
			fmt.Printf("Cache nicht lesbar %s: %v\n", r.Room, err)
			continue
		}
		// generate image
		httpResp, err := http.Post("http://172.20.0.4:72/generate", "application/json", bytes.NewBuffer(formatted))
		if err != nil {
			fmt.Printf("POST Fehler %s: %v\n", r.Room, err)
			continue
		}
		defer httpResp.Body.Close()
		_, err = io.ReadAll(httpResp.Body)
		if err != nil {
			fmt.Printf("Antwort lesen Fehler: %v\n", err)
			continue
		}

		// get image hex
		httpResp, err = http.Get("http://172.20.0.4:72/image?room=" + r.Room)
		if err != nil {
			fmt.Println("Fehler:", err)
			continue
		}
		defer httpResp.Body.Close()

		body, err := io.ReadAll(httpResp.Body) // ← GET Body lesen!
		if err != nil {
			fmt.Printf("Antwort lesen Fehler: %v\n", err)
			continue
		}

		fmt.Printf("Raw Antwort:", body)

		// JSON parsen
		var imgResp ImageResponse
		if err := json.Unmarshal(body, &imgResp); err != nil {
			fmt.Printf("JSON Parse Fehler: %v\n", err)
			continue
		}

		// Hex-String direkt speichern
		filename := filepath.Join(hexDir, r.Room+".hex")
		if err := os.WriteFile(filename, []byte(imgResp.Image), 0644); err != nil {
			fmt.Printf("Schreiben Fehler: %v\n", err)
			continue
		}
		fmt.Printf("Gespeichert: %s\n", filename)
	}

}

func now() time.Time {
	return time.Date(2026, 5, 11, 0, 0, 0, 0, time.Local)
}
