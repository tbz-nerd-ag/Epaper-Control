package untis

import (
	"encoding/json"
	"os"

	"github.com/robfig/cron/v3"
)

type Room struct {
	Room string `json:"room"`
}

func Schedule() {

	// pause 09:40 - 10:00
	// 11:30 - 11:45
	// 13:15 - 13:45
	// 15:15 - 15:30
	// 17:00 - 17:15

	c := cron.New()

	c.AddFunc("*/2 * * * *", func() {
		Get_room_from_json()
	})

	c.Start()
}

func Get_room_from_json() {
	data, err := os.ReadFile("untis/room.json")
	if err != nil {
		panic(err)
	}

	var rooms []Room
	if err := json.Unmarshal(data, &rooms); err != nil {
		panic(err)
	}

	for _, r := range rooms {
		Get_data(r.Room)
	}
}
