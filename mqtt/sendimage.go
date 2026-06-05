package mqtt

import (
	"Control/types"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func SendImage(c mqtt.Client, id string) {
	time.Sleep(5 * time.Second)
	//if wartung is enable send wartungs image
	if types.Config.Wartung {
		imageData, err := loadImage("mqtt/wartung.png")

		if err != nil {
			fmt.Printf("Fehler: %v\n", err)
			return
		}

		responseTopic := id + "/image"
		fmt.Println("Wartungsmodus ist aktiv")

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

		time.Sleep(5 * time.Second)
		sendsleep(c, id)
	}
}
