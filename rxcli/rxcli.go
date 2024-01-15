package rxcli

import (
	"github.com/NubeIO/rxclient"
)

var rxClient rxclient.RxClient

func RxClient() (rxclient.RxClient, error) {
	if rxClient == nil {
		r, err := rxclient.New()
		if err != nil {
			return nil, err
		}
		rxClient = r
		return rxClient, nil
	}
	return rxClient, nil
}
