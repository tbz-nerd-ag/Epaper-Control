package mqtt

import (
	"Control/handler"
	"Control/influx"
	"Control/types"
	"Control/untis"
	"fmt"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnecttoMQTT() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://172.20.0.5:1883")
	opts.SetClientID("Epaper-Control")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	subscribe(client)

	time.Sleep(22 * 1000)
	untis.Get_room_from_json()
	handler.PrepareJSON()
	handler.Getpicturehex()

	select {}
}

func subscribe(client mqtt.Client) {
	client.Subscribe("+/awake", 0, onAwake)
	client.Subscribe("+/gn", 0, onGN)
}

func onAwake(c mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) < 2 {
		return
	}
	id := parts[0]

	payload := strings.Split(string(msg.Payload()), ",")
	if len(payload) < 3 {
		return
	}
	batterypercent := payload[1]
	errorcode := payload[2]

	fmt.Printf("EPD %s ist wach, Akku %s%% und ErrorCode %s!\n", id, batterypercent, errorcode)

	battery, _ := strconv.Atoi(batterypercent)

	influx.SaveBatteryInflux(id, battery)

	room := types.GetRoomfromID(id)
	nightsleep, _ := types.GetNightsleep(room)
	if !nightsleep {
		SendImage(c, id)
	} else {
		sendsleep(c, id)
	}

}

func onGN(c mqtt.Client, msg mqtt.Message) {
	// Payload: "goodnight,600"
	parts := strings.Split(string(msg.Payload()), ",")
	if len(parts) < 2 {
		return
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		fmt.Printf("Fehler: %v\n", err)
		return
	}

	fmt.Printf("ESP32 schläft für %d Sekunden\n", seconds)
	wakeTime := time.Now().Add(time.Duration(seconds) * time.Second)
	fmt.Printf("Wacht auf um: %s\n", wakeTime.Format("15:04:05"))
}

func sendsleep(c mqtt.Client, id string) {
	responseTopic := id + "/sleep"

	var sekunden int

	//if wartun sleep for 30min if not sleep antil next lesson
	if !types.Config.Wartung {
		sekunden = handler.Getwakeuptime(types.GetRoomfromID(id)) * 60
	} else {
		sekunden = 30 * 60
	}

	send := c.Publish(responseTopic, 0, false, fmt.Sprintf("%d", sekunden))
	send.Wait()
	fmt.Println("EPD geht schlafen")

	//check if handler/cache/room.json is emtpy = night for epd
	room := types.GetRoomfromID(id)
	handler.IsNightSleep(id, room)
}
