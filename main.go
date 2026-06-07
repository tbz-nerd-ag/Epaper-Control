package main

import (
	"Control/handler"
	"Control/influx"
	"Control/mqtt"
	"Control/rest"
	"Control/types"
	"Control/untis"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	types.Loadconfig()
	types.Loadepd()
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

	r := gin.Default()

	auth := r.Group("/")
	auth.Use(rest.JWTMiddleware())
	{
		auth.GET("/get_wartung", rest.REST_GetWartung)
	}

	r.Run("0.0.0.0:80")

	mqtt.ConnecttoMQTT()
}
