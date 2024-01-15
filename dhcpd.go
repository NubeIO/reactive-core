package main

import (
	"github.com/NubeIO/reactive"
	"github.com/NubeIO/reactive-nodes/constants"
	dhcp "github.com/NubeIO/reactive-nodes/dhcpd"
	"github.com/NubeIO/rxlib"
	"github.com/gin-gonic/gin"
	"net/http"
)

var DHCP dhcpObject

type dhcpObject struct {
	rxlib.Object
	dhcp dhcp.DHCP
}

func NewDHCPObject(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	object := reactive.NewBaseObject(reactive.ObjectInfo(dhcpName, objectUUID, name, pluginName), bus)
	object.NewOutputPort(constants.Output, constants.Output, "bool")
	object.SetDetails(&rxlib.Details{
		Category:   categoryNetworkingDHCP,
		ObjectType: rxlib.Service,
	})
	object.AddDependencies(&rxlib.Dependencies{
		RequiresRouter: true,
	})
	object.AddObjectTypeTags("networking", "dhcp", "ip-address")
	return &dhcpObject{
		Object: object,
		dhcp:   dhcp.NewDHCP(""),
	}
}

func (n *dhcpObject) New(objectUUID, name string, bus *rxlib.EventBus, settings *rxlib.Settings) rxlib.Object {
	newObject := NewDHCPObject(objectUUID, name, bus, settings)
	return newObject
}

func (n *dhcpObject) getResp(c *gin.Context) {
	c.JSON(http.StatusOK, "fuck ya")
}

func (n *dhcpObject) NewRoute(r *gin.RouterGroup) {
	r.GET("myplugin/get", n.getResp)
}

func (n *dhcpObject) Start() {
	if n.NotLoaded() {

	}

}
