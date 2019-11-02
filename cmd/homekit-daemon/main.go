package main

import (
	"fmt"
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
	"homekit-daemon/pkg/accessories"
	"log"
	"strconv"
	"time"
)

func scaleInt(x, inMin, inMax, outMin, outMax int) int {
	return (x-inMin)*(outMax-outMin)/(inMax-inMin) + outMin
}

func getLastTemp(c client.Client) (result float64, err error) {
	q := client.NewQuery("SELECT \"value\" FROM \"temperature\" GROUP BY * ORDER BY DESC LIMIT 1", "sensor_data", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		result, err = strconv.ParseFloat(fmt.Sprintf("%s", response.Results[0].Series[0].Values[0][1]), 64)
	}
	return
}

func getLastHum(c client.Client) (result float64, err error) {
	q := client.NewQuery("SELECT \"value\" FROM \"humidity\" GROUP BY * ORDER BY DESC LIMIT 1", "sensor_data", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		result, err = strconv.ParseFloat(fmt.Sprintf("%s", response.Results[0].Series[0].Values[0][1]), 64)
	}
	return
}

func getLastAQ(c client.Client) (result int, err error) {
	q := client.NewQuery("SELECT \"quality\" FROM \"air_quality\" WHERE time > now() - 1m ORDER BY DESC LIMIT 1", "sensor_data", "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		rawInt, _ := strconv.ParseInt(fmt.Sprintf("%s", response.Results[0].Series[0].Values[0][1]), 10, 64)
		result = scaleInt(int(rawInt), 1, 100, 1, 5)
	}
	return
}

func main() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
	}
	defer c.Close()

	lastTemp, _ := getLastTemp(c)
	lastHum, _ := getLastHum(c)
	lastAQ, _ := getLastAQ(c)

	if err != nil {
		fmt.Println("Error while reading temp: ", err.Error())
	}

	//log.Debug.Enable()
	iaqSensorInfo := accessory.Info{
		Name: "IAQ Sensor",
	}
	acc := accessories.NewIAQSensor(
		iaqSensorInfo,
		lastHum,
		-0,
		100,
		0.1,
		lastTemp,
		-40,
		80,
		0.1,
		lastAQ,
		1,
		100,
		1,
	)

	config := hc.Config{Pin: "12344321", Port: "12345", StoragePath: "./db"}
	t, err := hc.NewIPTransport(config, acc.Accessory)

	if err != nil {
		log.Panic(err)
	}

	// Periodically toggle the switch's on characteristic
	go func() {
		for {
			lastTemp, _ := getLastTemp(c)
			lastHum, _ := getLastHum(c)
			lastAQ, _ := getLastAQ(c)

			if err != nil {
				fmt.Println("Error while reading temp: ", err.Error())
			}
			acc.Temperature.CurrentTemperature.SetValue(lastTemp)
			acc.Humidity.CurrentRelativeHumidity.SetValue(lastHum)
			acc.AirQuality.AirQuality.SetValue(lastAQ)

			time.Sleep(10 * time.Second)
		}
	}()

	hc.OnTermination(func() {
		<-t.Stop()
	})

	t.Start()
}
