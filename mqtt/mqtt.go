package mqtt

import (
	"fmt"
	"os"
	"path/filepath"
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

	fmt.Printf("EPD %s ist wach!\n", id)
	time.Sleep(5 * time.Second)
	wartung := false
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
		hexDir := "handler/image_hex"
		data, err := os.ReadFile(filepath.Join(hexDir, "2.105.hex"))
		if err != nil {
			fmt.Print("Hex nicht lesbar: %w", err)
		}

		// Hex-String → []byte
		parts := strings.Split(string(data), ", ")
		imageBytes := make([]byte, len(parts))
		for i, p := range parts {
			p = strings.TrimSpace(p)
			val, err := strconv.ParseUint(strings.TrimPrefix(p, "0x"), 16, 8)
			if err != nil {
				fmt.Print("Hex nicht lesbar: %w", err)
			}
			imageBytes[i] = byte(val)
		}
		fmt.Printf("Sende %d Bytes an %s/image\n", len(imageBytes), id)

		// Per MQTT senden
		topic := id + "/image"
		token := c.Publish(topic, 0, false, imageBytes)
		token.Wait()
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
	sekunden := 10 * 60

	send := c.Publish(responseTopic, 0, false, fmt.Sprintf("%d", sekunden))
	send.Wait()
	fmt.Println("EPD geht schlafen")
}
