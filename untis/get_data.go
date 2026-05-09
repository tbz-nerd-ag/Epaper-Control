package untis

import (
	"Control/types"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const outputDir = "untis/cache/"

func Get_room_from_json() {
	data, err := os.ReadFile("untis/room.json")
	if err != nil {
		panic(err)
	}

	var rooms []types.Room
	if err := json.Unmarshal(data, &rooms); err != nil {
		panic(err)
	}

	for _, r := range rooms {
		Get_data(r.Room)
	}
}

func Get_data(room string) {
	resp, err := http.Get("http://172.20.0.3:71/untis?room=" + room)
	if err != nil {
		fmt.Println("Untis Mircoservice nicht erreichbar,", err)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println("Json Decoding Err:", err)
		return
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Json MarshalIndent Err:", err)
		return
	}

	// Ordner erstellen falls nicht vorhanden
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Println("Ordner erstellen Err:", err)
		return
	}

	// Pfad zusammensetzen: "data/roomName.json"
	filename := filepath.Join(outputDir, room+".json")
	err = os.WriteFile(filename, formatted, 0644)
	if err != nil {
		fmt.Println("Datei schreiben Err:", err)
		return
	}
}
