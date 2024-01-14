package main

import (
	"fmt"
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"github.com/NubeIO/rxlib"
	"math/rand"
	"time"
)

// exports
var Trigger triggerFloat

// triggerFloat generates random values at regular intervals.
type triggerFloat struct {
	rxlib.Node
	stop chan struct{}
}

// NewTriggerNode creates a new triggerFloat with the given ID, name, EventBus, and Flow.
func NewTriggerNode(nodeUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Node {
	node := reactive.NewBaseNode(reactive.NodeInfo(trigger, nodeUUID, name, pluginName), bus)
	node.NewOutputPort(constants.Output, constants.Output, "float")
	return &triggerFloat{
		Node: node,
		stop: make(chan struct{}),
	}
}

func (n *triggerFloat) New(nodeUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Node {
	newNode := NewTriggerNode(nodeUUID, name, bus, settings)
	return newNode
}

func (n *triggerFloat) Start() {
	go func() {
		ticker := time.NewTicker(2000 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case <-n.stop:
					return // Stop triggering when the stopTrigger channel is closed
				default:
					ranValue := randFloat()
					out := &rxlib.Port{
						ID:        constants.Output,
						Name:      constants.Output,
						Value:     ranValue,
						Direction: "output",
						DataType:  "float",
					}
					fmt.Println("NEW VALUE:", "triggerFloat", out.Value)
					n.PublishMessage(out, true)
				}
			}
		}
	}()
}

func (n *triggerFloat) Delete() {
	close(n.stop)
	n.RemoveNodeFromRuntime()
}

func randFloat() float64 {
	rand.NewSource(time.Now().UnixNano())
	randomFloat := rand.Float64()*9 + 1
	return float64(int(randomFloat))
}
