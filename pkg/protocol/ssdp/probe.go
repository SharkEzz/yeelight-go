package ssdp

import (
	"bytes"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SharkEzz/yeelight-go/pkg/device"
	"github.com/SharkEzz/yeelight-go/pkg/protocol"
)

const MESSAGE_STRING = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb"

type YeelightProbe struct{}

func (y *YeelightProbe) Search() ([]*device.Device, error) {
	addr := &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1982}

	c, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return []*device.Device{}, err
	}
	defer c.Close()

	_, err = c.WriteToUDP([]byte(MESSAGE_STRING), addr)
	if err != nil {
		return []*device.Device{}, err
	}

	timeout := time.After(5 * time.Second)
	running := true

	lights := []*device.Device{}

	for running {
		select {
		case <-timeout:
			running = false
		default:
		}

		buf := make([]byte, 65536)
		c.SetReadDeadline(time.Now().Add(time.Second))
		c.Read(buf)
		buf = bytes.Trim(buf, "\x00")

		if len(buf) == 0 {
			// TODO: better handling
			continue
		}

		light := parseResponse(strings.Split(string(buf), "\r\n"))

		alreadyExist := false

		for _, storedLight := range lights {
			if !alreadyExist && storedLight.GetIP() == light.GetIP() {
				alreadyExist = true
			} else {
				continue
			}
		}

		if !alreadyExist {
			lights = append(lights, light)
		}
	}

	return lights, nil
}

func parseResponse(response []string) *device.Device {
	location := regexp.MustCompile(`(?i)yeelight://([0-9.]+):([0-9]+)`).FindAllStringSubmatch(response[4], 1)
	ip := location[0][1]
	portStr := location[0][2]

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		// TODO: better handling
		panic(err)
	}

	light := device.NewDevice(protocol.NewYeelightClient(ip, uint16(port)))

	return light
}
