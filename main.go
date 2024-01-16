package main

import (
	"fmt"
	pprint "github.com/NubeIO/reactive-nodes/helpers/print"
	"github.com/NubeIO/reactive-nodes/rxcli"
	"github.com/NubeIO/reactive/plugins"
)

var Plugin pluginExport

const pluginName = "my plugin"
const pluginVersion = "v1.0"

const categoryNetworkingDHCP = "networking-dhcp"
const dhcpName = "dhcp"

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

	go rxcli.RxClient()

	var err error
	e.AddCategory(categoryTime)
	//e.AddCategory(categoryCount)
	//e.AddCategory(categoryModbus)
	//err = e.AddObject(categoryTime, count, countExport)
	//if err != nil {
	//	fmt.Println(err)
	//}
	err = e.AddObject(categoryTime, trigger, triggerExport)
	if err != nil {
		fmt.Println(err)
	}

	// modbus
	//err = e.AddObject(categoryModbus, modbusNetworkName, modbusNetworkExport)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//err = e.AddChildObject(categoryModbus, modbusNetworkName, modbusDeviceName, modbusDeviceExport)
	//fmt.Println(err)
	//err = e.AddChildObject(categoryModbus, modbusDeviceName, modbusPointName, modbusPointExport)
	//fmt.Println(err)

	pprint.PrintJOSN(e)
	fmt.Println(err)
	return e
}
