package main

import (
	"fmt"
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	"github.com/NubeIO/rxlib"
	"github.com/NubeIO/schema"
)

var Count countObject

// countObject represents a object that counts incoming messages and sends out the count value.
type countObject struct {
	//*reactive.BaseObject
	rxlib.Object
	count int
}

// NewCountObject creates a new countObject with the given ID, name, and EventBus.
func NewCountObject(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	object := reactive.NewBaseObject(reactive.ObjectInfo(count, objectUUID, name, pluginName), bus)
	object.NewInputPort(constants.Input, constants.Input, "any")
	object.NewOutputPort(constants.Output, constants.Output, "float")
	object.SetHotFix()
	return &countObject{
		Object: object,
		count:  0,
	}
}

func (n *countObject) New(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	newObject := NewCountObject(objectUUID, name, bus, settings)
	newObject.AddSchema()
	return newObject
}

func (n *countObject) Start() {
	if n.NotLoaded() {
		n.SetLoaded(true)
		inputChannel, exists := n.BusChannel(constants.Input)

		if !exists {
			fmt.Printf("Input channel for target input %s does not exist\n", constants.Input)
			return
		}
		for {
			_, ok := <-inputChannel
			if !ok {
				return
			}
			// Increment the count for each incoming message
			n.count++
			// Create a Port with the count value and send it to the output
			countPort := &rxlib.Port{
				ID:        constants.Output,
				Name:      constants.Output,
				Value:     float64(n.count), // Convert count to float64
				Direction: "output",
				DataType:  "float", // Data type of the output
			}
			n.PublishMessage(countPort, true)
		}
	}

}

type countObjectSettings struct {
	StartCount int
}

func (n *countObject) AddSchema() {

	//builder := schema.NewSchemaBuilder("HEY")
	//
	//builder.NewString("exampleString", "Example String", true, 3, 10, "default value")
	//builder.NewNumber("exampleNumber", "Example Num", true, nil, nil, 11)
	//
	//ui := schema.UI{}
	//ui.AddUIOrder([]string{"exampleString", "exampleNumber"})
	//out := &schema.Generated{
	//	schema: builder.Build(),
	//	UI:     ui,
	//}

	builder := schema.NewSchemaBuilder("IF/THEN").
		SetProperty("holidayType", schema.Property{Type: "string", Enum: []string{"snow", "beach"}})

	newCon := make(map[string]schema.Property)
	aa := 8.0
	newCon["a"] = schema.NewNumber("want a new aa", false, nil, &aa, 1)
	newCon["b"] = schema.NewNumber("want a new bb", false, nil, &aa, 1)

	builder.AddCondition(schema.ConditionalStructure{
		If: schema.Condition{
			Properties: map[string]schema.Property{"holidayType": {Const: "beach"}},
		},
		Then: schema.Condition{
			Properties: newCon,
		},
	})

	newCon2 := make(map[string]schema.Property)

	newCon2["a"] = schema.NewString("want a new aa str", false, 1, 100, "hey")
	newCon2["bb"] = newCon["b"]

	// Dynamic condition for 'snow' holiday
	builder.AddCondition(schema.ConditionalStructure{
		If: schema.Condition{
			Properties: map[string]schema.Property{"holidayType": {Const: "snow"}},
		},
		Then: schema.Condition{
			Properties: newCon2,
		},
	})

	out := &schema.Generated{
		Schema: builder.Build(),
	}
	n.BuildSchema(out)
}
