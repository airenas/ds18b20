// Copyright 2016 yryz Author. All Rights Reserved.

package ds18b20

import (
	"errors"
	"strconv"
	"strings"
)

type reader interface {
	Get(name string) (string, error)
}

//Sensor to provide data
type Sensor struct {
	dReader reader

	retryCount int // 0 - no retry
	key        string
}

//NewSensor initializes sensor object
func NewSensor(key string) (*Sensor, error) {
	return NewSensorR(key, 0)
}

//NewSensorR initializes sensor object with retry
func NewSensorR(key string, retryOnError int) (*Sensor, error) {
	return newSensor(key, retryOnError, &fileReader{})
}

func newSensor(key string, retryOnError int, dReader reader) (*Sensor, error) {
	if key == "" {
		return nil, errors.New("No key provided")
	}
	if retryOnError < 0 {
		return nil, errors.New("Fail check: retryOnError >= 0")
	}
	if dReader == nil {
		return nil, errors.New("No reader provided")
	}
	var res Sensor
	res.key = key
	res.retryCount = retryOnError
	res.dReader = dReader
	return &res, nil
}

//Temperature reads sensor temperature
func (s *Sensor) Temperature() (float64, error) {
	t, err := s.read()
	for rc := 0; rc < s.retryCount && err != nil; rc++ {
		t, err = s.read()
	}
	return t, err
}

func (s *Sensor) read() (float64, error) {
	d, err := s.dReader.Get(s.key)
	if err != nil {
		return 0.0, err
	}
	return extract(d)
}

func extract(raw string) (float64, error) {
	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, ErrReadSensor
	}

	td := strings.TrimSpace(raw[i+2 : len(raw)-1])
	if td == "85000" { // some sensor error
		return 0.0, ErrReadSensor
	}
	c, err := strconv.ParseFloat(td, 64)
	if err != nil {
		return 0.0, ErrReadSensor
	}
	return c / 1000.0, nil
}
