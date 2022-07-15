package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/SharkEzz/yeelight-go/pkg/bulb"
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

	go func() {
		for light.IsConnected() {
			res := <-light.ResponseChannel

			fmt.Printf("%+v\n", res)
		}
	}()

	light.GetProp([]any{"power", "bright"})
	time.Sleep(time.Second * 2)
}
