package watchdog

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type (
	ChannelMode12 byte
	ChannelMode3  byte

	Params struct {
		Alarm            time.Duration
		ResetPress       time.Duration
		HardResetHold    time.Duration
		HardResetRelease time.Duration
		HardResetPress   time.Duration
		Channel1         ChannelMode12
		Channel2         ChannelMode12
		ResetLimit       byte
		Channel3         ChannelMode3
		TempThreshold    byte
	}
)

const (
	ChannelMode12Off   ChannelMode12 = 0
	ChannelMode12Reset ChannelMode12 = 1
	ChannelMode12Power ChannelMode12 = 2
	ChannelMode12Open  ChannelMode12 = 3
	ChannelMode12Close ChannelMode12 = 4

	ChannelMode3Off    ChannelMode3 = 0
	ChannelMode3Input  ChannelMode3 = 1
	ChannelMode3Output ChannelMode3 = 2
	ChannelMode3Temp   ChannelMode3 = 3
)

func (m ChannelMode12) String() string {
	switch m {
	case ChannelMode12Off:
		return "OFF"
	case ChannelMode12Reset:
		return "RESET"
	case ChannelMode12Power:
		return "POWER"
	case ChannelMode12Open:
		return "OPEN"
	case ChannelMode12Close:
		return "CLOSE"
	default:
		return ""
	}
}

func (m ChannelMode3) String() string {
	switch m {
	case ChannelMode3Off:
		return "OFF"
	case ChannelMode3Input:
		return "INPUT"
	case ChannelMode3Output:
		return "OUTPUT"
	case ChannelMode3Temp:
		return "TEMP"
	default:
		return ""
	}
}

func decodeNumber(p []byte) byte {
	x, err := strconv.ParseInt(string(p), 16, 8)
	if err != nil {
		return 0
	}
	return byte(x)
}

func decodeDuration(p []byte) time.Duration {
	return time.Duration(decodeNumber(p))
}

func decodeMode12(p []byte) ChannelMode12 {
	return ChannelMode12(decodeNumber(p))
}

func decodeMode3(p []byte) ChannelMode3 {
	return ChannelMode3(decodeNumber(p))
}

func encodeNumber(p, size byte) []byte {
	leadZero := p < 16 && size > 1
	x := strconv.FormatInt(int64(p), 16)
	if leadZero {
		x = "0" + x
	}
	return []byte(strings.ToUpper(x))
}

func encodeDuration(p, base time.Duration) []byte {
	return encodeNumber(byte(p/base), 1)
}

func encodeMode12(p ChannelMode12) []byte {
	return encodeNumber(byte(p), 1)
}

func encodeMode3(p ChannelMode3) []byte {
	return encodeNumber(byte(p), 1)
}

func (p *Params) decode(raw []byte) {
	p.Alarm = decodeDuration(raw[0:1]) * time.Minute
	p.ResetPress = decodeDuration(raw[1:2]) * 100 * time.Millisecond
	p.HardResetHold = decodeDuration(raw[2:3]) * time.Second
	p.HardResetRelease = decodeDuration(raw[3:4]) * time.Second
	p.HardResetPress = decodeDuration(raw[4:5]) * 100 * time.Millisecond
	p.Channel1 = decodeMode12(raw[5:6])
	p.Channel2 = decodeMode12(raw[6:7])
	p.Channel3 = decodeMode3(raw[8:9])
	p.ResetLimit = decodeNumber(raw[7:8])
	p.TempThreshold = decodeNumber(raw[9:11])
}

func (p *Params) encode() []byte {
	raw := make([]byte, 0, 11)
	raw = append(raw, encodeDuration(p.Alarm, time.Minute)...)
	raw = append(raw, encodeDuration(p.ResetPress, 100*time.Millisecond)...)
	raw = append(raw, encodeDuration(p.HardResetHold, time.Second)...)
	raw = append(raw, encodeDuration(p.HardResetRelease, time.Second)...)
	raw = append(raw, encodeDuration(p.HardResetPress, 100*time.Millisecond)...)
	raw = append(raw, encodeMode12(p.Channel1)...)
	raw = append(raw, encodeMode12(p.Channel2)...)
	raw = append(raw, encodeNumber(p.ResetLimit, 1)...)
	raw = append(raw, encodeMode3(p.Channel3)...)
	raw = append(raw, encodeNumber(p.TempThreshold, 2)...)
	return raw
}

func (p *Params) String() (s string) {
	s += fmt.Sprintf("Alarm:            %v\n", p.Alarm)
	s += fmt.Sprintf("ResetPress:       %v\n", p.ResetPress)
	s += fmt.Sprintf("HardResetHold:    %v\n", p.HardResetHold)
	s += fmt.Sprintf("HardResetRelease: %v\n", p.HardResetRelease)
	s += fmt.Sprintf("HardResetPress:   %v\n", p.HardResetPress)
	s += fmt.Sprintf("Channel 1:        %v\n", p.Channel1)
	s += fmt.Sprintf("Channel 2:        %v\n", p.Channel2)
	s += fmt.Sprintf("Channel 3:        %v\n", p.Channel3)
	s += fmt.Sprintf("ResetLimit:       %v\n", p.ResetLimit)
	s += fmt.Sprintf("TempThreshold:    %v\n", p.TempThreshold)
	return
}
