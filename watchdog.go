package watchdog

import (
	"fmt"
	"sync"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

type (
	command string
	answer  string
)

const (
	WatchDogVID = "0483"
	WatchDogPID = "A26D"

	cmdPing        command = "~U"
	cmdVersion     command = "~I"
	cmdReadParams  command = "~F"
	cmdWriteParams command = "~W"
	cmdControl     command = "~M"
	cmdLight       command = "~L"
	cmdTest        command = "~T"

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
		if port.VID == WatchDogVID && port.PID == WatchDogPID {
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

func (dog *WatchDog) send(input []byte) (output []byte, err error) {

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

	time.Sleep(time.Second)

	output = make([]byte, 32)
	nRead, err := port.Read(output)
	if err != nil || nRead == 0 {
		return nil, err
	}

	return output[:nRead], nil
}

func (dog *WatchDog) sendCommand(cmd command, input []byte) (output []byte, err error) {
	cmdInput := []byte(cmd)
	cmdInput = append(cmdInput, input...)

	cmdOutput, err := dog.send(cmdInput)
	if err != nil {
		return
	}

	if string(cmdOutput[:2]) != string(cmd) {
		err = fmt.Errorf("wrong answer")
		return
	}

	output = cmdOutput[2:]
	return
}

func (dog *WatchDog) Ping() (bool, error) {
	output, err := dog.send([]byte(cmdPing))
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

func (dog *WatchDog) Control1() error {
	_, err := dog.sendCommand(cmdControl, []byte{1})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) Control2() error {
	_, err := dog.sendCommand(cmdControl, []byte{2})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) LightOff() error {
	_, err := dog.sendCommand(cmdLight, []byte{0})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) LightOn() error {
	_, err := dog.sendCommand(cmdLight, []byte{1})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) Reset() error {
	_, err := dog.sendCommand(cmdTest, []byte{1})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) HardReset() error {
	_, err := dog.sendCommand(cmdTest, []byte{2})
	if err != nil {
		return err
	}
	return nil
}

func (dog *WatchDog) PowerOff() error {
	_, err := dog.sendCommand(cmdTest, []byte{3})
	if err != nil {
		return err
	}
	return nil
}
