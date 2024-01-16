package main

import (
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	dhcp "github.com/NubeIO/reactive-nodes/dhcpd"
	"github.com/NubeIO/rxlib"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

var DHCP dhcpObject

type dhcpObject struct {
	rxlib.Object
	dhcp     dhcp.DHCP
	filePath string
}

func NewDHCPObject(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	object := reactive.NewBaseObject(reactive.ObjectInfo(dhcpName, objectUUID, name, pluginName), bus)
	object.NewOutputPort(constants.Output, constants.Output, "bool")
	object.SetDetails(&rxlib.Details{
		Category:   categoryNetworkingDHCP,
		ObjectType: rxlib.Service,
	})
	object.AddObjectTypeRequirement(rxlib.RequirementMaxOne())
	object.AddObjectTypeRequirement(rxlib.RequirementWebRouter())
	object.AddObjectTypeTags(rxlib.Networking, rxlib.IpAddress)
	var filePath = ""
	return &dhcpObject{
		Object:   object,
		dhcp:     dhcp.NewDHCP(filePath),
		filePath: filePath,
	}
}

func (n *dhcpObject) New(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	newObject := NewDHCPObject(objectUUID, name, bus, settings)
	return newObject
}

func (n *dhcpObject) Start() {
	if n.NotLoaded() {

	}
	if n.HaltFlag() {

	}
}

var fileNotFoundHaltKey = "fileNotFound"

func (n *dhcpObject) RunValidation() {
	found := fileExists(n.filePath)
	if found {
		n.NewHalt("fileNotFound", "dhcpd file was not found", "this os type could be incorrect")
	}
}

func (n *dhcpObject) resetHalt() {
	found := fileExists(n.filePath)
	if found {
		n.DeleteValidation(fileNotFoundHaltKey)
	}
}

func (n *dhcpObject) NewRoute(r *gin.RouterGroup) {
	r.GET("myplugin/get", n.getResp)
}

func (n *dhcpObject) getResp(c *gin.Context) {
	c.JSON(http.StatusOK, "fuck ya")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
