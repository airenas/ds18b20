// Copyright 2016 yryz Author. All Rights Reserved.

package ds18b20

import (
	"io/ioutil"
)

type fileReader struct {
}

func (fr *fileReader) Get(sensor string) (string, error) {
	d, err := ioutil.ReadFile("/sys/bus/w1/devices/" + sensor + "/w1_slave")
	if err != nil {
		return "", err
	}
	return string(d), nil
}
