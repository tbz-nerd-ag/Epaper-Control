package handler

import (
	"encoding/json"
	"os"
	"time"
)

func Getwakeuptime(room string) int {
	data, err := os.ReadFile("handler/cache/" + room + ".json")
	if err != nil {
		return 10 * 60
	}
	var resp Response
	if err := json.Unmarshal(data, &resp); err != nil {
		return 10 * 60
	}
	if len(resp.Lessons) == 0 {
		return 10 * 60
	}

	now := time.Now()
	nowMinutes := now.Hour()*60 + now.Minute()

	nextlesson := resp.Lessons[0]
	startMinutes := parseTime(nextlesson.StartTime)

	if startMinutes > nowMinutes {
		return startMinutes
	}
	return parseTime(nextlesson.EndTime)
}

func parseTime(t string) int {
	if len(t) != 4 {
		return 10 * 60
	}
	hours := int(t[0]-'0')*10 + int(t[1]-'0')
	mins := int(t[2]-'0')*10 + int(t[3]-'0')
	return hours*60 + mins
}
