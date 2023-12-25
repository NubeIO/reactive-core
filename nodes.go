package main

import (
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"github.com/NubeIO/reactive/schemas"
	"math/rand"
	"time"
)

// triggerNode generates random values at regular intervals.
type triggerNode struct {
	*reactive.BaseNode
	stop chan struct{}
}

// NewTriggerNode creates a new triggerNode with the given ID, name, EventBus, and Flow.
func NewTriggerNode(nodeUUID, name string, bus *reactive.EventBus, opts *reactive.NodeOpts) reactive.Node {
	node := reactive.NewBaseNode("trigger-float", nodeUUID, name, bus)
	node.NewOutputPort(constants.Output, constants.Output, "float")
	if opts != nil {
		node.Meta = opts.Meta
		node.AddToNodesMap(nodeUUID, node)
	}
	return &triggerNode{BaseNode: node}
}

func (n *triggerNode) New(nodeUUID, name string, bus *reactive.EventBus, opts *reactive.NodeOpts) reactive.Node {
	newNode := NewTriggerNode(nodeUUID, name, bus, opts)
	//Node[newNode.GetID()] = newNode
	return newNode
}

func (n *triggerNode) Start() {
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

func (n *triggerNode) Delete() {
	close(n.stop)
	n.RemoveFromNodesMap()
}

func (n *triggerNode) BuildSchema() *schemas.Schema {
	return nil
}

func randFloat() float64 {
	rand.NewSource(time.Now().UnixNano())
	randomFloat := rand.Float64()*9 + 1
	return float64(int(randomFloat))
}

type triggerBool struct {
	*reactive.BaseNode
	stop chan struct{}
}

func NewTriggerBool(nodeUUID, name string, bus *reactive.EventBus, opts *reactive.NodeOpts) reactive.Node {
	node := reactive.NewBaseNode("trigger-bool", nodeUUID, name, bus)
	node.NewOutputPort(constants.Output, constants.Output, "float")
	if opts != nil {
		node.Meta = opts.Meta
		node.AddToNodesMap(nodeUUID, node)
	}
	return &triggerBool{BaseNode: node}
}

func (n *triggerBool) New(nodeUUID, name string, bus *reactive.EventBus, opts *reactive.NodeOpts) reactive.Node {
	newNode := NewTriggerBool(nodeUUID, name, bus, opts)
	//Node[newNode.GetID()] = newNode
	return newNode
}

func (n *triggerBool) Start() {
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

func (n *triggerBool) Delete() {
	close(n.stop)
	n.RemoveFromNodesMap()
}

func (n *triggerBool) BuildSchema() *schemas.Schema {
	return nil
}

// exports
var TriggerNode triggerNode
var TriggerBool triggerBool
