package main

import (
	"fmt"
	"time"

	"github.com/qovan/watchdog"
)

func main() {
	for _, name := range watchdog.EnumerateWatchDog() {
		dog := watchdog.NewWatchDog(name)
		version, _ := dog.GetVersion()
		fmt.Println(name, " : ", version)
		params, _ := dog.ReadParams()
		fmt.Println(params)
		var err error
		print := func(err error) {
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("OK")
			}
		}
		for range 10 {
			fmt.Print("PING: ")
			_, err = dog.Ping()
			print(err)

			switch params.Channel1 {
			case watchdog.ChannelMode12Reset:

				fmt.Print("TOUCH 1: ")
				err = dog.Touch(watchdog.Channel1)
				print(err)

				fmt.Print("RESET: ")
				err = dog.Reset()
				print(err)

			case watchdog.ChannelMode12Open, watchdog.ChannelMode12Close:

				fmt.Print("TURN ON: ")
				err = dog.TurnON(watchdog.Channel1)
				print(err)

				fmt.Print("TURN OFF: ")
				err = dog.TurnOFF(watchdog.Channel1)
				print(err)
			}

			switch params.Channel2 {
			case watchdog.ChannelMode12Power:

				fmt.Print("CONTROL 2: ")
				err = dog.Touch(watchdog.Channel2)
				print(err)

				fmt.Print("HARD RESET: ")
				err = dog.HardReset()
				print(err)

				fmt.Print("POWER OFF: ")
				err = dog.PowerOff()
				print(err)

			case watchdog.ChannelMode12Open, watchdog.ChannelMode12Close:

				fmt.Print("TURN ON: ")
				err = dog.TurnON(watchdog.Channel2)
				print(err)

				fmt.Print("TURN OFF: ")
				err = dog.TurnOFF(watchdog.Channel2)
				print(err)

			}

			fmt.Print("PAUSE ON: ")
			err = dog.Pause(watchdog.ON)
			print(err)

			time.Sleep(time.Second)

			fmt.Print("PAUSE OFF: ")
			err = dog.Pause(watchdog.OFF)
			print(err)
		}
	}
}
