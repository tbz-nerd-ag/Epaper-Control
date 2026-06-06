package handler

import (
	"Control/types"
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
		//no lesson in the room for the day long sleep is required or a wish from mr. scheichel
		nightsleep, _ := types.GetNightsleep(room)
		if nightsleep {
			data, err := os.ReadFile("untis/cache/" + room + ".json")
			if err != nil {
				return 10 * 60
			}
			var resp Response
			if err := json.Unmarshal(data, &resp); err != nil {
				return 10 * 60
			}

			nextlesson := resp.Lessons[0]
			nextLessonDate, _ := time.ParseInLocation("2006-01-02", nextlesson.Date, time.Local)

			sleeptime := parseTime(nextlesson.StartTime)

			lessonTime := nextLessonDate.Add(time.Duration(sleeptime) * time.Minute)
			sleepDuration := time.Until(lessonTime)
			sleeptimeMinutes := int(sleepDuration.Minutes())

			if sleepDuration > 6*time.Hour {
				return 6 * 60
			}
			return sleeptimeMinutes
		}
	}

	now := time.Now()
	nowMinutes := now.Hour()*60 + now.Minute()

	nextlesson := resp.Lessons[0]
	startMinutes := parseTime(nextlesson.StartTime)

	//if lesson start time is behind -> endtime
	if startMinutes > nowMinutes {
		return startMinutes - nowMinutes
	}
	endMinutes := parseTime(nextlesson.EndTime)
	return endMinutes - nowMinutes
}

func parseTime(t string) int {
	if len(t) != 4 {
		return 10 * 60
	}
	hours := int(t[0]-'0')*10 + int(t[1]-'0')
	mins := int(t[2]-'0')*10 + int(t[3]-'0')
	return hours*60 + mins
}
