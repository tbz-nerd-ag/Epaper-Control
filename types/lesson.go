package types

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
