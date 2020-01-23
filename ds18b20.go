// Copyright 2016 yryz Author. All Rights Reserved.

package ds18b20

import (
	"errors"
	"io/ioutil"
	"strings"
)

//ErrReadSensor is a sensor read error
var ErrReadSensor = errors.New("failed to read sensor temperature")

// SensorIDs get all connected sensor IDs as array
func SensorIDs() ([]string, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves")
	if err != nil {
		return nil, err
	}

	sensors := strings.Split(string(data), "\n")
	if len(sensors) > 0 {
		sensors = sensors[:len(sensors)-1]
	}
	return sensors, nil
}
