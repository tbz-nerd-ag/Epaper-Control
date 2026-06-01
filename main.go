package main

import (
	"Control/handler"
	"Control/influx"
	"Control/mqtt"
	"Control/types"
	"Control/untis"
	"fmt"

	"github.com/robfig/cron/v3"
)

func main() {
	types.Loadconfig()
	types.Loadepd()

	// pause 09:40 - 10:00
	// 11:30 - 11:45
	// 13:15 - 13:45
	// 15:15 - 15:30
	// 17:00 - 17:15
	c := cron.New()

	for _, task := range types.Config.Task_time_cron {
		task := task
		fmt.Print(task + " registiert ...")
		c.AddFunc(task, func() { // task direkt, nicht task.Cron
			fmt.Println("Cron wird ausgeführt:", task)
			untis.Get_room_from_json()
			handler.PrepareJSON()
			handler.Getpicturehex()
		})
	}
	c.Start()

	influx.InitInflux()
	mqtt.ConnecttoMQTT()
}
