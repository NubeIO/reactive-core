package main

import (
	"encoding/json"
	"fmt"
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"github.com/NubeIO/reactive-nodes/helpers/pointers"
	pprint "github.com/NubeIO/reactive-nodes/helpers/print"
	"github.com/grid-x/modbus"
	"log"
	"time"
)

var ModbusNetwork modbusNetwork
var ModbusDevice modbusDevice
var ModbusPoint modbusPoint

type modbusNetwork struct {
	*reactive.BaseNode
	pollInterval time.Duration // Interval between polls
	stopChannel  chan struct{} // Channel to signal stopping of polling
	client       modbus.Client
	tcpClient    *modbus.TCPClientHandler
	rtuClient    *modbus.RTUClientHandler
	isRTUNetwork bool
}

func NewModbusNetwork(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings, opts *reactive.Options) reactive.Node {
	node := reactive.NewBaseNode(modbusNetworkName, nodeUUID, name, bus, opts)
	node.NewInputPort(constants.Input, constants.Input, "any")
	node.NewOutputPort(constants.Output, constants.Output, "float")
	node.SetDetails(&reactive.Details{
		Category: categoryModbus,
	})
	return &modbusNetwork{
		BaseNode:     node,
		pollInterval: time.Second * 2,
	}
}

func (n *modbusNetwork) New(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings, opts *reactive.Options) reactive.Node {
	newNode := NewModbusNetwork(nodeUUID, name, bus, settings, opts)
	newNode.AddSchema()
	return newNode
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
				fmt.Println("data", err)
				return
			}

			pprint.PrintJOSN(mb)
			//modbusPoint.

			//parsedPoint, ok := device.(*modbusPoint)
			//if !ok {
			//	continue
			//}

			data, err := n.client.ReadCoils(1, 1) // Example usage
			fmt.Println(data)
			if err != nil {
				log.Printf("Error reading Modbus device %s: %v", device.GetID(), err)
				continue
			} else {
				log.Printf("modbus-read %s", data)

			}
			// Update point value
			device.SetLastValueChildNode(point.GetUUID(), &reactive.Port{
				ID:    constants.Output,
				Value: data,
			})
		}
	}
}

type modbusDevice struct {
	*reactive.BaseNode
	deviceAddr int
}

func NewModbusDevice(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings, opts *reactive.Options) reactive.Node {
	node := reactive.NewBaseNode(modbusDeviceName, nodeUUID, name, bus, opts)
	node.NewInputPort(constants.Input, constants.Input, "any")
	node.NewOutputPort(constants.Output, constants.Output, "float")
	node.SetDetails(&reactive.Details{
		Category: categoryModbus,
		ParentID: pointers.NewString(modbusNetworkName),
	})
	return &modbusDevice{
		BaseNode:   node,
		deviceAddr: 1,
	}
}

func (n *modbusDevice) New(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings, opts *reactive.Options) reactive.Node {
	newNode := NewModbusDevice(nodeUUID, name, bus, settings, opts)
	newNode.AddSchema()
	return newNode
}

type modbusPoint struct {
	*reactive.BaseNode
	*pointSettings
}

func NewModbusPoint(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings, opts *reactive.Options) reactive.Node {
	node := reactive.NewBaseNode(modbusPointName, nodeUUID, name, bus, opts)
	node.NewInputPort(constants.Input, constants.Input, "any")
	node.NewOutputPort(constants.Output, constants.Output, "float")
	node.SetDetails(&reactive.Details{
		Category: categoryModbus,
		ParentID: pointers.NewString(modbusDeviceName),
	})
	n := &modbusPoint{
		BaseNode:      node,
		pointSettings: nil,
	}
	n.AddSettings(settings)
	return n
}

func (n *modbusPoint) New(nodeUUID, name string, bus *reactive.EventBus, settings *reactive.Settings, opts *reactive.Options) reactive.Node {
	newNode := NewModbusPoint(nodeUUID, name, bus, settings, opts)
	newNode.AddSchema()
	return newNode
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

func (n *modbusPoint) AddSettings(settings *reactive.Settings) {
	out := &pointSettings{
		Register: 3,
		Function: coil,
		Request:  read,
	}
	fmt.Println("############################")
	marshal, err := json.Marshal(settings)
	if err != nil {
		fmt.Println(11111, err, "MODBUS POINT SETTINGS")

		return
	}
	err = json.Unmarshal(marshal, &settings)
	n.AddData(modbusPointName, out)
	n.pointSettings = out
}
