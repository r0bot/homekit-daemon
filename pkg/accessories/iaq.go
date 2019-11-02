package accessories

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type IAQSensor struct {
	*accessory.Accessory

	Temperature *service.TemperatureSensor
	Humidity    *service.HumiditySensor
	AirQuality  *service.AirQualitySensor
}

// NewTemperatureSensor returns a Thermometer which implements model.Thermometer.
func NewIAQSensor(info accessory.Info, hum, humMin, humMax, humSteps, temp, tempMin, tempMax, tempSteps float64, aq, aqMin, aqMax, aqStep int) *IAQSensor {
	acc := IAQSensor{}
	acc.Accessory = accessory.New(info, accessory.TypeSensor)

	//Temperature
	acc.Temperature = service.NewTemperatureSensor()
	acc.Temperature.CurrentTemperature.SetValue(temp)
	acc.Temperature.CurrentTemperature.SetMinValue(tempMin)
	acc.Temperature.CurrentTemperature.SetMaxValue(tempMax)
	acc.Temperature.CurrentTemperature.SetStepValue(tempSteps)
	acc.AddService(acc.Temperature.Service)

	//Humidity
	acc.Humidity = service.NewHumiditySensor()
	acc.Humidity.CurrentRelativeHumidity.SetValue(hum)
	acc.Humidity.CurrentRelativeHumidity.SetMinValue(humMin)
	acc.Humidity.CurrentRelativeHumidity.SetMaxValue(humMax)
	acc.Humidity.CurrentRelativeHumidity.SetStepValue(humSteps)
	acc.AddService(acc.Humidity.Service)

	//AirQuality
	acc.AirQuality = service.NewAirQualitySensor()
	acc.AirQuality.AirQuality.SetValue(aq)
	acc.AirQuality.AirQuality.SetMinValue(aqMin)
	acc.AirQuality.AirQuality.SetMaxValue(aqMax)
	acc.AirQuality.AirQuality.SetStepValue(aqStep)
	acc.AddService(acc.AirQuality.Service)

	return &acc
}
