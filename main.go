package main

import (
	"errors"
	"github.com/Depau/go-iio-sensor-proxy"
	"github.com/godbus/dbus"
	"log"
	"time"
)

func SwayRotate(orientation string) (err error) {
	log.Println("Rotating:", orientation)

	var degrees string

	switch orientation {
	case "normal":
		degrees = "0"
	case "right-up":
		degrees = "90"
	case "bottom-up":
		degrees = "180"
	case "left-up":
		degrees = "270"
	default:
		return errors.New("Unrecognized orientation: " + orientation)
	}

	var outputs []string
	outputs, err = GetOutputNames()

	if err != nil {
		return
	}

	for _, output := range outputs {
		err = SwayMsg(nil, "output", output, "transform", degrees)

		if err != nil {
			return
		}
	}

	// Inputs can't be rotated yet

	//var inputs []SwayInput
	//
	//for _, input := range inputs {
	//	if input.Type == "touchpad" || (input.Type == "pointer" && strings.Contains(input.Name, "TrackPoint")) {
	//		err = SwayMsg(nil, "input", "eDP-1", "transform", degrees)
	//	}
	//}

	return
}

func Claim(sensorProxy sensorproxy.SensorProxy) {
	if err := sensorProxy.ClaimAccelerometer(); err != nil {
		log.Fatal("Failed to claim accelerometer:", err)
	}
}

func Release(sensorProxy sensorproxy.SensorProxy) {
	if err := sensorProxy.ReleaseAccelerometer(); err != nil {
		log.Fatal("Failed to claim accelerometer:", err)
	}
}

func getOrientationAndRotate(sensorProxy sensorproxy.SensorProxy, previous string) (newOrientation string) {
	var err error
	newOrientation, err = sensorProxy.GetAccelerometerOrientation()

	if err != nil {
		log.Fatal("Failed to get orientation:", err)
	}

	if newOrientation != previous {
		err = SwayRotate(newOrientation)

		if err != nil {
			log.Fatal("Unable to rotate:", err)
		}
	}

	return
}

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatal("Failed to connect to system bus:", err)
	}

	sensorProxy, err := sensorproxy.NewSensorProxyFromBus(conn)
	if err != nil {
		log.Fatal("Failed to get sensorProxy object from DBus:", err)
	}

	hasAccel, err := sensorProxy.HasAccelerometer()
	if err != nil {
		log.Fatal("Failed to check whether device has an accelerometer:", err)
	}
	if !hasAccel {
		log.Fatal("No accelerometer found")
	}

	//noinspection GoInfiniteFor
	currentOrientation := "undefined"

	Claim(sensorProxy)
	defer Release(sensorProxy)

	for {
		currentOrientation = getOrientationAndRotate(sensorProxy, currentOrientation)
		time.Sleep(1e9)
	}
}
