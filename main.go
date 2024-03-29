package main

import (
	"os"

	"github.com/nicholasjackson/consul-smi-controller/consul"
	"github.com/servicemeshinterface/smi-controller-sdk/sdk"
	"github.com/servicemeshinterface/smi-controller-sdk/sdk/controller"
)

func main() {
	api, err := consul.New(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		panic(err)
	}

	// register our lifecycle callbacks with the controller
	sdk.API().RegisterV1Alpha(api)

	// create and start a the controller
	config := controller.DefaultConfig()
	controller.Start(config)
}
