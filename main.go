package main

import (
	"fmt"
	"github.com/NubeIO/reactive/plugins"
)

var Plugin pluginExport

const pluginName = "my plugin"
const pluginVersion = "v1.0"

const categoryTime = "time"
const trigger = "trigger"
const triggerExport = "Trigger"

const categoryCount = "count"
const count = "count"
const countExport = "Count"

const categoryModbus = "modbus"
const modbusNetworkName = "modbus-network"
const modbusNetworkExport = "ModbusNetwork"
const modbusDeviceName = "modbus-device"
const modbusDeviceExport = "ModbusDevice"
const modbusPointName = "modbus-point"
const modbusPointExport = "ModbusPoint"

type pluginExport struct{}

func (p *pluginExport) Get() *plugins.Export {
	e := plugins.NewPlugin(pluginName, pluginVersion, "a new plugin")

	e.AddCategory(categoryTime)
	e.AddCategory(categoryCount)
	e.AddCategory(categoryModbus)
	e.AddNode(categoryTime, count, countExport)
	e.AddNode(categoryTime, trigger, triggerExport)

	// modbus
	e.AddNode(categoryModbus, modbusNetworkName, modbusNetworkExport)
	err := e.AddChildNode(categoryModbus, modbusNetworkName, modbusDeviceName, modbusDeviceExport)
	err = e.AddChildNode(categoryModbus, modbusDeviceName, modbusPointName, modbusPointExport)
	fmt.Println(err)
	return e
}
