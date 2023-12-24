package main

import (
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"github.com/NubeIO/reactive/schemas"
	"math/rand"
	"time"
)

// TriggerNode generates random values at regular intervals.
type TriggerNode struct {
	*reactive.BaseNode
	stop chan struct{}
}

// NewTriggerNode creates a new TriggerNode with the given ID, name, EventBus, and Flow.
func NewTriggerNode(nodeUUID, name string, bus *reactive.EventBus, opts *reactive.NodeOpts) reactive.Node {
	node := reactive.NewBaseNode("trigger-float", nodeUUID, name, bus)
	node.NewOutputPort(constants.Output, constants.Output, "float")
	if opts != nil {
		node.Meta = opts.Meta
		node.AddToNodesMap(nodeUUID, node)
	}
	return &TriggerNode{BaseNode: node}
}

func (n *TriggerNode) New(nodeUUID, name string, bus *reactive.EventBus, opts *reactive.NodeOpts) reactive.Node {
	return NewTriggerNode(nodeUUID, name, bus, opts)
}

func (n *TriggerNode) Start() {
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
					out := &reactive.Port{
						ID:        constants.Output,
						Name:      constants.Output,
						Value:     ranValue,
						Direction: "output",
						DataType:  "float",
					}
					n.PublishMessage(out, true)
				}
			}
		}
	}()
}

func (n *TriggerNode) Delete() {
	close(n.stop)
	n.RemoveFromNodesMap()
}

func (n *TriggerNode) BuildSchema() *schemas.Schema {
	return nil
}

func randFloat() float64 {
	rand.NewSource(time.Now().UnixNano())
	randomFloat := rand.Float64()*9 + 1
	return float64(int(randomFloat))
}

var Node TriggerNode
