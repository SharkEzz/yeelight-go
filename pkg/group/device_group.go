package group

import (
	"fmt"
	"sync"

	"github.com/SharkEzz/yeelight-go/pkg/device"
)

// A DeviceGroup hold multiple devices and allow to send commands to all of them.
//
// The devices added to the group must NOT not be connected to the real device.
type DeviceGroup struct {
	devices []*device.Device
}

// NewDeviceGroup create a new DeviceGroup.
func NewDeviceGroup(devices []*device.Device) *DeviceGroup {
	return &DeviceGroup{
		devices,
	}
}

// Add a device to the group.
func (d *DeviceGroup) Add(device *device.Device) error {
	if device.IsConnected() {
		return fmt.Errorf("device is already connected")
	}

	for _, item := range d.devices {
		if item.IP == device.IP {
			return fmt.Errorf("device is already in the group")
		}
	}

	d.devices = append(d.devices, device)

	return nil
}

// Connect all the devices in the group.
func (d *DeviceGroup) Connect() error {
	for _, device := range d.devices {
		err := device.Connect()
		if err != nil {
			return err
		}
	}

	return nil
}

// Disconnect all the devices in the group.
func (d *DeviceGroup) Disconnect() {
	for _, device := range d.devices {
		device.Disconnect()
	}
}

// Toggle all the devices in the group.
func (d *DeviceGroup) Toggle() {
	wg := sync.WaitGroup{}
	wg.Add(len(d.devices))

	for _, d := range d.devices {
		go func(d *device.Device) {
			d.Toggle()
			wg.Done()
		}(d)
	}

	wg.Wait()
}
