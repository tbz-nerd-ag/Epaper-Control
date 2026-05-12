package main

import (
	"Control/handler"
	"Control/mqtt"
	"Control/types"
	"Control/untis"
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	types.Loadconfig()

	// pause 09:40 - 10:00
	// 11:30 - 11:45
	// 13:15 - 13:45
	// 15:15 - 15:30
	// 17:00 - 17:15

	c := cron.New()

	c.AddFunc(types.Config.Task_time_cron, func() {
		fmt.Println("Cron wird ausgeführt")
		untis.Get_room_from_json()
		handler.PrepareJSON()
		handler.Getpicturehex()
	})

	c.Start()
	mqtt.ConnecttoMQTT()
	//untis.Get_room_from_json()
}
