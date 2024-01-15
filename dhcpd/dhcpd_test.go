package dhcp

import (
	"fmt"
	"testing"
)

func TestNewDHCP(t *testing.T) {
	file := "dhcpcd.conf"
	got := NewDHCP(file)
	var err error
	ok := got.FileExists()
	fmt.Println(ok)

	err = got.SetFaceAsDHCPOrRemove("eth1")
	fmt.Println(err)

	//err = got.SetFaceAsStatic("eth0", "192.168.15.10", "255.0.0.0", "192.168.1.1")
	fmt.Println(err)

}
