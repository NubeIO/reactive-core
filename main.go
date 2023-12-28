package main

import (
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

type pluginExport struct{}

func (p *pluginExport) Get() *plugins.Export {
	e := plugins.NewPlugin(pluginName, pluginVersion, "a new plugin")

	e.AddCategory(categoryTime)
	e.AddNode(categoryTime, count, countExport)
	e.AddNode(categoryTime, trigger, triggerExport)

	return e
}
