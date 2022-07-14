package main

import (
	"flag"
	"fmt"

	"github.com/SharkEzz/yeelight-go/pkg/bulb"
	"github.com/SharkEzz/yeelight-go/pkg/bulb/command"
)

func main() {
	name := flag.String("name", "default", "The light name")
	ip := flag.String("ip", "", "The light IP address (in IPv4 format)")

	flag.Parse()

	if *ip == "" {
		panic(fmt.Errorf("ip must be set"))
	}

	light := bulb.NewBulb(*name, *ip)
	err := light.Connect()
	if err != nil {
		panic(err)
	}
	defer light.Disconnect()

	str, err := light.SendCommand(command.SetBright(100, "smooth", 1000))
	if err != nil {
		panic(err)
	}
	fmt.Print(str)
}
