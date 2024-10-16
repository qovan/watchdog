package watchdog

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDecodeNumber(t *testing.T) {
	data := []byte("0123456789ABCDEF")
	dataErr := []byte("GHIXYZ&%$")

	for i := range data {
		want := byte(i)
		value := decodeNumber(data[i : i+1])
		assert.Equal(t, want, value)
	}

	for i := range dataErr {
		want := byte(0)
		value := decodeNumber(dataErr[i : i+1])
		assert.Equal(t, want, value)
	}
}

func TestDecodeDuration(t *testing.T) {
	data := []byte("0123456789ABCDEF")
	dataErr := []byte("GHIXYZ&%$")

	for i := range data {
		want := time.Duration(i)
		value := decodeDuration(data[i : i+1])
		assert.Equal(t, want, value)
	}

	for i := range dataErr {
		want := time.Duration(0)
		value := decodeDuration(dataErr[i : i+1])
		assert.Equal(t, want, value)
	}
}

func TestDecodeMode12(t *testing.T) {
	data := []byte("01234")
	dataErr := []byte("GHIXYZ&%$")

	for i := range data {
		want := ChannelMode12(i)
		value := decodeMode12(data[i : i+1])
		assert.Equal(t, want, value)
	}

	for i := range dataErr {
		want := ChannelMode12(0)
		value := decodeMode12(dataErr[i : i+1])
		assert.Equal(t, want, value)
	}
}

func TestDecodeMode3(t *testing.T) {
	data := []byte("0123")
	dataErr := []byte("GHIXYZ&%$")

	for i := range data {
		want := ChannelMode3(i)
		value := decodeMode3(data[i : i+1])
		assert.Equal(t, want, value)
	}

	for i := range dataErr {
		want := ChannelMode3(0)
		value := decodeMode3(dataErr[i : i+1])
		assert.Equal(t, want, value)
	}
}

func TestDecodeParams(t *testing.T) {
	wantValue := Params{
		Alarm:            time.Duration(15) * time.Minute,
		ResetPress:       time.Duration(200) * time.Millisecond,
		HardResetHold:    time.Duration(6) * time.Second,
		HardResetRelease: time.Duration(2) * time.Second,
		HardResetPress:   time.Duration(200) * time.Millisecond,
		Channel1:         ChannelMode12Reset,
		Channel2:         ChannelMode12Power,
		ResetLimit:       10,
		Channel3:         ChannelMode3Off,
		TempThreshold:    80,
	}

	data := []byte("F262212A050")

	value := Params{}
	value.decode(data)

	assert.Equal(t, wantValue, value)
}

func TestEncodeNumber(t *testing.T) {
	data := []byte("0123456789ABCDEF")

	for i, wantValue := range data {
		value := encodeNumber(byte(i), 1)
		assert.Equal(t, []byte{wantValue}, value)
	}
}

func TestEncodeDuration(t *testing.T) {
	data := []byte("0123456789ABCDEF")

	for i, wantValue := range data {
		value := encodeDuration(time.Duration(i)*time.Second, time.Second)
		assert.Equal(t, []byte{wantValue}, value)
	}
}

func TestEncodeMode12(t *testing.T) {
	data := []byte("01234")

	for i, wantValue := range data {
		value := encodeMode12(ChannelMode12(i))
		assert.Equal(t, []byte{wantValue}, value)
	}
}

func TestEncodeMode3(t *testing.T) {
	data := []byte("0123")

	for i, wantValue := range data {
		value := encodeMode3(ChannelMode3(i))
		assert.Equal(t, []byte{wantValue}, value)
	}
}

func TestEncodeParams(t *testing.T) {
	wantValue := []byte("F262212A050")
	params := Params{
		Alarm:            time.Duration(15) * time.Minute,
		ResetPress:       time.Duration(200) * time.Millisecond,
		HardResetHold:    time.Duration(6) * time.Second,
		HardResetRelease: time.Duration(2) * time.Second,
		HardResetPress:   time.Duration(200) * time.Millisecond,
		Channel1:         ChannelMode12Reset,
		Channel2:         ChannelMode12Power,
		ResetLimit:       10,
		Channel3:         ChannelMode3Off,
		TempThreshold:    80,
	}

	value := params.encode()

	assert.Equal(t, wantValue, value)
}

func TestEncodeDecodeParams(t *testing.T) {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	for i := range 16 {
		wantValue := Params{
			Alarm:            time.Duration((values[0]+i)%16) * time.Minute,
			ResetPress:       time.Duration((values[1]+i)%16) * 100 * time.Millisecond,
			HardResetHold:    time.Duration((values[2]+i)%16) * time.Second,
			HardResetRelease: time.Duration((values[3]+i)%16) * time.Second,
			HardResetPress:   time.Duration((values[4]+i)%16) * 100 * time.Millisecond,
			Channel1:         ChannelMode12((values[5] + i) % 5),
			Channel2:         ChannelMode12((values[6] + i) % 5),
			ResetLimit:       byte((values[7] + i) % 16),
			Channel3:         ChannelMode3((values[8] + i) % 4),
			TempThreshold:    byte((values[9] + i) % 255),
		}

		value := Params{}
		value.decode(wantValue.encode())

		assert.Equal(t, wantValue, value)
	}
}
