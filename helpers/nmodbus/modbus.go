package nmodbus

type ModbusNetworks struct {
	UUID    string
	Devices map[string]*ModbusDevices
}

type ModbusDevices struct {
	Points     map[string]ModbusPoint
	UUID       string
	DeviceIP   string
	DevicePort int
	DeviceAddr int
}

type ModbusPoint struct {
	UUID     string
	Register int
	Function string // e.g., "coil"
	Request  string // e.g., "read" or "write"
}

// NewModbusNetwork creates a new Modbus network.
func NewModbusNetwork(uuid string) *ModbusNetworks {
	return &ModbusNetworks{
		UUID:    uuid,
		Devices: make(map[string]*ModbusDevices),
	}
}

// AddDevice adds a new device to the network.
func (m *ModbusNetworks) AddDevice(device *ModbusDevices) {
	m.Devices[device.UUID] = device
}

func (m *ModbusNetworks) GetDevices() map[string]*ModbusDevices {
	return m.Devices
}

// NewDevice creates a new Modbus device.
func NewDevice(uuid, deviceIP string, devicePort, deviceAddr int) *ModbusDevices {
	return &ModbusDevices{
		Points:     make(map[string]ModbusPoint),
		UUID:       uuid,
		DeviceIP:   deviceIP,
		DevicePort: devicePort,
		DeviceAddr: deviceAddr,
	}
}

// NewPoint creates a new Modbus point.
func NewPoint(uuid string, register int, function, request string) ModbusPoint {
	return ModbusPoint{
		UUID:     uuid,
		Register: register,
		Function: function,
		Request:  request,
	}
}

// GetPoints returns all points of a device.
func (d *ModbusDevices) GetPoints() map[string]ModbusPoint {
	return d.Points
}

// AddPoint adds a new point to the device.
func (d *ModbusDevices) AddPoint(point ModbusPoint) {
	d.Points[point.UUID] = point
}

// GetPoint returns a point with the specified UUID.
func (d *ModbusDevices) GetPoint(uuid string) (ModbusPoint, bool) {
	point, exists := d.Points[uuid]
	return point, exists
}
