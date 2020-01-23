package ds18b20

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dsplit      = ":::"
	correctData = "58 01 55 00 7f ff 0c 10 ff : crc=ff YES\n58 01 55 00 7f ff 0c 10 ff t=21500"
	wrongData85 = "58 01 55 00 7f ff 0c 10 ff : crc=ff YES\n58 01 55 00 7f ff 0c 10 ff t=85000"
)

func TestExtract_FailsOnNoData(t *testing.T) {
	_, err := extract("")
	assert.NotNil(t, err)
}

func TestExtract_FailsOnWrongData(t *testing.T) {
	_, err := extract("xxxx\nxxx")
	assert.NotNil(t, err)
	_, err = extract("xxxx\nxxx t=aaaa")
	assert.NotNil(t, err)
}

func TestExtract_ReturnsValue(t *testing.T) {
	f, err := extract(correctData)
	assert.Nil(t, err)
	assert.Equal(t, 21.5, f)

	f, err = extract("58 01 55 00 7f ff 0c 10 ff : crc=ff YES\n58 01 55 00 7f ff 0c 10 ff t=-21500")
	assert.Nil(t, err)
	assert.Equal(t, -21.5, f)
}

func TestExtract_FailsOn85(t *testing.T) {
	_, err := extract(wrongData85)
	assert.NotNil(t, err)
}

func TestExtract_FailsOnCrc(t *testing.T) {
	_, err := extract("58 01 55 00 7f ff 0c 10 ff : crc=ff NO\n58 01 55 00 7f ff 0c 10 ff t=25000")
	assert.NotNil(t, err)
}

func TestSensor_FailsOnNoData(t *testing.T) {
	_, err := NewSensor("")
	assert.NotNil(t, err)
	_, err = NewSensorR("", 0)
	assert.NotNil(t, err)
}

func TestSensor_FailsOnWrongData(t *testing.T) {
	_, err := NewSensorR("ss", -1)
	assert.NotNil(t, err)
}

func TestSensor_ReturnsValue(t *testing.T) {
	s, err := newSensor("ss", 0, newTReader(correctData))
	assert.Nil(t, err)
	f, err := s.Temperature()
	assert.Nil(t, err)
	assert.Equal(t, 21.5, f)
}

func TestSensor_ReturnFailure(t *testing.T) {
	s, err := newSensor("ss", 0, newTReader(wrongData85))
	assert.Nil(t, err)
	_, err = s.Temperature()
	assert.NotNil(t, err)
}

func TestSensor_FailsOnError(t *testing.T) {
	s, err := newSensor("ss", 0, newTReader(""))
	assert.Nil(t, err)
	_, err = s.Temperature()
	assert.NotNil(t, err)
}

func TestSensor_Retries(t *testing.T) {
	s, err := newSensor("ss", 2, newTReader(wrongData85+dsplit+wrongData85+dsplit+correctData))
	assert.Nil(t, err)
	f, err := s.Temperature()
	assert.Nil(t, err)
	assert.Equal(t, 21.5, f)
}

func TestSensor_RetriesJustOneTime(t *testing.T) {
	rd := newTReader(wrongData85 + dsplit + wrongData85 + dsplit + correctData)
	s, err := newSensor("ss", 1, rd)
	assert.Nil(t, err)
	_, err = s.Temperature()
	assert.NotNil(t, err)
	assert.NotNil(t, 2, rd.i)
}

type tReader struct {
	d []string
	i int
}

func newTReader(data string) *tReader {
	return &tReader{d: strings.Split(data, dsplit)}
}

func (r *tReader) Get(sensor string) (string, error) {
	if len(r.d) < r.i {
		return "", errors.New("No data")
	}
	i := r.i
	r.i++
	return r.d[i], nil
}
