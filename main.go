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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Control/docs"
)

// @title 	DoorSign Control MircoService API
// @version v1.0
// @host localhost:80
// @BasePath        /
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
		auth.GET("/get_wartung_sleep", rest.REST_GetWartungSleepTime)

		auth.POST("/post_wartung", rest.REST_PostWartung)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go mqtt.ConnecttoMQTT()

	r.Run("0.0.0.0:80")
}
