package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"github.com/NubeIO/reactive-nodes/helpers/pointers"
	"github.com/NubeIO/reactive-nodes/rxcli"
	"github.com/NubeIO/rxclient"
	"github.com/NubeIO/rxlib"
	"github.com/grid-x/modbus"
	"time"
)

var ModbusNetwork modbusNetwork
var ModbusDevice modbusDevice
var ModbusPoint modbusPoint

type modbusNetwork struct {
	rxlib.Object
	pollInterval time.Duration // Interval between polls
	stopChannel  chan struct{} // Channel to signal stopping of polling
	client       modbus.Client
	tcpClient    *modbus.TCPClientHandler
	rtuClient    *modbus.RTUClientHandler
	isRTUNetwork bool
}

func NewModbusNetwork(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	object := reactive.NewBaseObject(reactive.ObjectInfo(modbusNetworkName, objectUUID, name, pluginName), bus)
	object.NewInputPort(constants.Input, constants.Input, "any")
	object.NewOutputPort(constants.Output, constants.Output, "float")
	object.AddDefinedChildObjects(modbusDeviceName)
	object.SetDetails(&rxlib.Details{
		Category:   categoryModbus,
		ObjectType: rxlib.Driver,
	})
	n := &modbusNetwork{
		Object:       object,
		pollInterval: time.Second * 2,
	}
	return n
}

func (n *modbusNetwork) New(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	newObject := NewModbusNetwork(objectUUID, name, bus, settings)
	return newObject
}

func (n *modbusNetwork) setClient() modbus.Client {
	n.tcpClient = modbus.NewTCPClientHandler("localhost:10502")
	client := modbus.NewClient(n.tcpClient)
	n.tcpClient.Connect()
	defer n.tcpClient.Close()
	return client
}

func (n *modbusNetwork) setDeviceAddr(addr int) {
	n.tcpClient.SetSlave(byte(addr))
}

// Start the polling process
func (n *modbusNetwork) Start() {
	if n.Loaded() {
		return
	}
	n.client = n.setClient()
	go func() {
		ticker := time.NewTicker(n.pollInterval)
		for {
			select {
			case <-ticker.C:
				n.pollDevices()
			case <-n.stopChannel:
				ticker.Stop()
				return
			}
		}
	}()
	n.SetLoaded(true)
}

// pollDevices performs the Modbus read operation for each point in each device
func (n *modbusNetwork) pollDevices() {

	devices := n.GetChildsByType(modbusDeviceName)
	for _, device := range devices {
		parsedDevice, ok := device.(*modbusDevice)
		if !ok {
			continue
		}
		n.setDeviceAddr(parsedDevice.deviceAddr)
		points := device.GetChildsByType(modbusPointName)

		for _, point := range points {
			mb := &pointSettings{}
			err := point.GetDataByKey(modbusPointName, &mb)
			if err != nil {
				//logger.Debug("read", "func.GetDataByKey()", "err:", err.Error())
				return
			}
			fmt.Println("read-coil", "addr:", 1, "count:", 1)
			//logger.Info("read-coil", "addr:", 1, "count:", 1)
			data, err := n.client.ReadCoils(1, 1) // Example usage
			if err != nil {
				fmt.Println("read-coil", "addr:", 1, "count:", 1, "err:", err.Error())
				//logger.Error()
				continue
			} else {

			}
			// Update point value
			device.SetLastValueChildObject(point.GetUUID(), &rxlib.Port{
				ID:    constants.Output,
				Value: data,
			})
		}
	}
}

type modbusDevice struct {
	rxlib.Object
	deviceAddr int
}

func NewModbusDevice(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	object := reactive.NewBaseObject(reactive.ObjectInfo(modbusDeviceName, objectUUID, name, pluginName), bus)
	object.NewInputPort(constants.Input, constants.Input, "any")
	object.NewOutputPort(constants.Output, constants.Output, "float")
	object.AddDefinedChildObjects(modbusPointName)
	object.SetDetails(&rxlib.Details{
		Category:   categoryModbus,
		ObjectType: rxlib.Driver,
		ParentID:   pointers.NewString(modbusNetworkName),
	})
	return &modbusDevice{
		Object:     object,
		deviceAddr: 1,
	}
}

func (n *modbusDevice) New(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	newObject := NewModbusDevice(objectUUID, name, bus, settings)
	return newObject
}

type modbusPoint struct {
	rxlib.Object
	*pointSettings
	rxClient rxclient.RxClient
}

func NewModbusPoint(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	object := reactive.NewBaseObject(reactive.ObjectInfo(modbusPointName, objectUUID, name, pluginName), bus)
	object.NewInputPort(constants.Input, constants.Input, "any")
	object.NewOutputPort(constants.Output, constants.Output, "float")
	object.SetDetails(&rxlib.Details{
		Category:   categoryModbus,
		ObjectType: rxlib.Driver,
		ParentID:   pointers.NewString(modbusDeviceName),
	})
	rx, err := rxcli.RxClient()
	if err != nil {
		object.AddValidationResult("rx-client-init", fmt.Sprintf("error on init: %v", err))
	}
	n := &modbusPoint{
		Object:        object,
		pointSettings: nil,
		rxClient:      rx,
	}
	n.AddSettings(settings)
	return n
}

func (n *modbusPoint) New(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	newObject := NewModbusPoint(objectUUID, name, bus, settings)
	return newObject
}

// RunValidation example of a validation on adding a new point
func (n *modbusPoint) RunValidation() {
	validation := make(map[string]any)

	validation["addingPoint"] = "the same register type as already been added before"
	n.SetValidationResult(validation)
}

type functionType string
type requestType string

const (
	coil functionType = "coil"
)
const (
	read requestType = "read"
)

type pointSettings struct {
	Register uint16       `json:"register"`
	Function functionType `json:"function"` // e.g., "coil"
	Request  requestType  `json:"request"`  // e.g., "read" or "write"
}

func (n *pointSettings) register() uint16 {
	return n.Register
}
func (n *pointSettings) function() functionType {
	return n.Function
}

func (n *pointSettings) request() requestType {
	return n.Request
}

func (n *modbusPoint) AddSettings(settings *rxlib.Settings) {
	out := &pointSettings{
		Register: 3,
		Function: coil,
		Request:  read,
	}

	marshal, err := json.Marshal(settings)
	if err != nil {

		return
	}
	err = json.Unmarshal(marshal, &settings)
	n.AddData(modbusPointName, out)
	n.pointSettings = out
}
