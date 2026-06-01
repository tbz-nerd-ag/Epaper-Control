package influx

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

var influxClient influxdb2.Client

func InitInflux() {
	influxClient = influxdb2.NewClient(
		"http://172.20.0.6:8086",
		"jdhuezr7384uind3ih8rufh3in°!",
	)
}

func SaveBatteryInflux(id string, battery int) {
	writeAPI := influxClient.WriteAPI(
		"epaper",
		"battery",
	)

	p := influxdb2.NewPointWithMeasurement("battery").AddTag("device", id).AddField("percent", battery).SetTime(time.Now())

	writeAPI.WritePoint(p)
	writeAPI.Flush()
}
