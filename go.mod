module github.com/NubeIO/reactive-nodes

go 1.21.5

replace github.com/NubeIO/rxlib => /home/aidan/code/go/rxlib

replace github.com/NubeIO/reactive => /home/aidan/code/go/reactive

require (
	github.com/NubeIO/reactive v0.0.7
	github.com/NubeIO/rxlib v0.0.2
	github.com/NubeIO/schema v0.0.1
	github.com/grid-x/modbus v0.0.0-20230713135356-d9fefd3ae5a5
)

require (
	github.com/NubeIO/rxclient v0.0.2 // indirect
	github.com/NubeIO/unixclient v0.0.2 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grid-x/serial v0.0.0-20191104121038-e24bc9bf6f08 // indirect
	golang.org/x/net v0.19.0 // indirect
)
