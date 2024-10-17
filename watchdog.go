package watchdog

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type (
	Channel rune
	State   rune
	command string
	answer  string
)

const (
	WatchDogVID = "0483"
	WatchDogPID = "a26d"

	Channel1 Channel = '1'
	Channel2 Channel = '2'
	Channel3 Channel = '3'

	OFF State = '0'
	ON  State = '1'

	cmdPing        command = "~U"
	cmdTurnON      command = "~S"
	cmdTurnOFF     command = "~R"
	cmdTest        command = "~T"
	cmdPause       command = "~P"
	cmdTouch       command = "~M"
	cmdLight       command = "~L"
	cmdBootloader  command = "~D"
	cmdVersion     command = "~I"
	cmdInput       command = "~G"
	cmdWriteParams command = "~W"
	cmdReadParams  command = "~F"

	ansPing answer = "~A"
)

type (
	WatchDog struct {
		name  string
		mode  serial.Mode
		mutex sync.Mutex
	}
)

func EnumerateWatchDog() (names []string) {
	names = make([]string, 0, 1)
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return
	}
	for _, port := range ports {
		if strings.ToLower(port.VID) == WatchDogVID && strings.ToLower(port.PID) == WatchDogPID {
			names = append(names, port.Name)
		}
	}
	return
}

func NewWatchDog(name string) *WatchDog {
	return &WatchDog{
		name: name,
		mode: serial.Mode{
			BaudRate: 9600,
		},
		mutex: sync.Mutex{},
	}
}

func (dog *WatchDog) send(input []byte, doRead bool) (output []byte, err error) {

	dog.mutex.Lock()
	defer dog.mutex.Unlock()

	port, err := serial.Open(dog.name, &dog.mode)
	if err != nil {
		return nil, err
	}
	defer port.Close()

	nWrite, err := port.Write(input)
	if err != nil || nWrite == 0 {
		return nil, err
	}

	if doRead {
		time.Sleep(time.Second)

		output = make([]byte, 32)
		nRead, err := port.Read(output)
		if err != nil || nRead == 0 {
			return nil, err
		}

		return output[:nRead], nil
	} else {
		return nil, nil
	}
}

func (dog *WatchDog) sendCommand(cmd command, input []byte) (output []byte, err error) {
	cmdInput := []byte(cmd)
	cmdInput = append(cmdInput, input...)

	hasOutput := true
	switch cmd {
	case cmdTest:
		hasOutput = false
	}

	cmdOutput, err := dog.send(cmdInput, hasOutput)
	if err != nil {
		return
	}

	if hasOutput {
		if cmdOutput != nil {
			if string(cmdOutput[:2]) != string(cmd) {
				err = fmt.Errorf("wrong answer: %s", string(cmdOutput))
				return
			}

			output = cmdOutput[2:]
		}
	}

	return
}

func (dog *WatchDog) Ping() (bool, error) {
	output, err := dog.send([]byte(cmdPing), true)
	if err != nil {
		return false, err
	}

	if string(ansPing) != string(output[:len(ansPing)]) {
		return false, fmt.Errorf("wrong ping answer %s", string(output[:len(ansPing)]))
	}

	return true, nil
}

func (dog *WatchDog) GetName() string {
	return dog.name
}

func (dog *WatchDog) GetVersion() (string, error) {
	output, err := dog.sendCommand(cmdVersion, []byte{})
	if len(output) > 0 && output[len(output)-1] == 10 {
		output = output[:len(output)-1]
	}
	return string(output), err
}

func (dog *WatchDog) ReadParams() (*Params, error) {
	output, err := dog.sendCommand(cmdReadParams, []byte{})
	if err != nil {
		return nil, err
	}
	params := Params{}
	params.decode(output)
	return &params, nil
}

func (dog *WatchDog) WriteParams(params *Params) error {
	raw := params.encode()
	_, err := dog.sendCommand(cmdWriteParams, raw)
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) Touch(ch Channel) error {
	switch ch {
	case Channel1, Channel2:
		_, err := dog.sendCommand(cmdTouch, []byte{byte(ch)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dog *WatchDog) Light(st State) error {
	_, err := dog.sendCommand(cmdLight, []byte{byte(st)})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) Reset() error {
	_, err := dog.sendCommand(cmdTest, []byte{'1'})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) HardReset() error {
	_, err := dog.sendCommand(cmdTest, []byte{'2'})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) PowerOff() error {
	_, err := dog.sendCommand(cmdTest, []byte{'3'})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) TurnON(ch Channel) error {
	switch ch {
	case Channel1, Channel2:
		_, err := dog.sendCommand(cmdTurnON, []byte{byte(ch)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dog *WatchDog) TurnOFF(ch Channel) error {
	switch ch {
	case Channel1, Channel2:
		_, err := dog.sendCommand(cmdTurnOFF, []byte{byte(ch)})
		if err != nil {
			return err
		}
	}
	return nil
}

func (dog *WatchDog) Pause(st State) error {
	_, err := dog.sendCommand(cmdPause, []byte{byte(st)})
	if err != nil {
		return err
	}
	return nil
}
