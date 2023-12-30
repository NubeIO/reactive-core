package nmodbus

import (
	pprint "github.com/NubeIO/reactive-nodes/helpers/print"
	"testing"
)

func TestNewModbusNetwork(t *testing.T) {

	net := NewModbusNetwork("net-1")

	dev := NewDevice("dev-1", "192", 502, 1)
	net.AddDevice(dev)

	pnt := NewPoint("dev-1", 1, "coil", "read")
	dev.AddPoint(pnt)
	pprint.PrintJOSN(net)

}
