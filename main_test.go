package main

import (
	pprint "github.com/NubeIO/reactive-nodes/helpers/print"
	"testing"
)

func Test_instance_Get(t *testing.T) {
	p := &pluginExport{}

	nodes := p.Get()

	pprint.PrintJOSN(nodes)
}
