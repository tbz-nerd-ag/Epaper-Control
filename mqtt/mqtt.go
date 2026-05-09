package mqtt

import (
	"fmt"
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

	select {}
}

func subscribe(client mqtt.Client) {
	client.Subscribe("+/awake", 0, onAwake)
}

func onAwake(c mqtt.Client, msg mqtt.Message) {
	parts := strings.Split(msg.Topic(), "/")
	if len(parts) < 2 {
		return
	}
	id := parts[0]

	fmt.Printf("EPD %s ist wach!\n", id)
	time.Sleep(5 * time.Second)
	wartung := true
	if wartung {
		imageData, err := loadImage("mqtt/wartung.png")

		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}

		responseTopic := id + "/image"

		token := c.Publish(responseTopic, 0, false, imageData)
		fmt.Printf("Bild gesendet an %s\n", responseTopic)
		token.Wait()

		time.Sleep(5 * time.Second)
		sendsleep(c, id)
	} else {

	}
}

func sendsleep(c mqtt.Client, id string) {
	responseTopic := id + "/sleep"
	sekunden := 10 * 60

	send := c.Publish(responseTopic, 0, false, sekunden)
	send.Wait()
	fmt.Println("EPD geht schlafen")
}
