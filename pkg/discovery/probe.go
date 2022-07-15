package discovery

import (
	"bytes"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/SharkEzz/yeelight-go/pkg/device"
)

const MESSAGE_STRING string = "M-SEARCH * HTTP/1.1\r\nHOST: 239.255.255.250:1982\r\nMAN: \"ssdp:discover\"\r\nST: wifi_bulb"

type Probe struct {
	Bulbs []*device.Device
}

func NewProbe() *Probe {
	return &Probe{}
}

func (p *Probe) Search() ([]*device.Device, error) {
	addr := &net.UDPAddr{IP: net.IPv4(239, 255, 255, 250), Port: 1982}

	c, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	_, err = c.WriteToUDP([]byte(MESSAGE_STRING), addr)
	if err != nil {
		return nil, err
	}

	// Allow 2 seconds for all the bulbs to respond.
	timeout := time.After(time.Second * 2)
	running := true

	bulbs := []*device.Device{}

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
			running = false
			break
		}

		bulb := parseResponse(strings.Split(string(buf), "\r\n"))

		alreadyExist := false

		for _, storedBulb := range bulbs {
			if !alreadyExist && storedBulb.IP == bulb.IP {
				alreadyExist = true
			} else {
				continue
			}
		}

		if !alreadyExist {
			bulbs = append(bulbs, bulb)
		}
	}

	return bulbs, nil
}

func parseResponse(response []string) *device.Device {
	location := regexp.MustCompile(`(?i)yeelight://([0-9.]+):([0-9]+)`).FindAllStringSubmatch(response[4], 1)
	ip := location[0][1]

	light := device.NewDevice("", ip)

	return light
}
