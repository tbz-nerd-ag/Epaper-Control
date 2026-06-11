package influx

import (
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func SaveRefreshTimeInflux(id string, d time.Duration) {
	writeAPI := influxClient.WriteAPI(
		"epaper",
		"display",
	)
	p := influxdb2.NewPointWithMeasurement("refreshtime").
		AddTag("device", id).
		AddField("seconds", d.Seconds()).
		SetTime(time.Now())
	writeAPI.WritePoint(p)
	writeAPI.Flush()
	fmt.Printf("InfluxDB: Refreshtime für %s gespeichert: %v\n", id, d.Seconds())
}
