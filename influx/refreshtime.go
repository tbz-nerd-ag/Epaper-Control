package influx

import (
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func SaveRefreshTimeInflux(id string, d time.Duration) {
	writeAPI := influxClient.WriteAPI(
		"epaper",
		"refreshtime",
	)
	p := influxdb2.NewPointWithMeasurement("refreshtime").
		AddTag("device", id).
		AddField("seconds", d.Seconds()).
		SetTime(time.Now())
	writeAPI.WritePoint(p)
	writeAPI.Flush()
}
