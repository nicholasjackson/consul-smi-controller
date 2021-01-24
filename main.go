package main

import (
	"github.com/nicholasjackson/consul-smi-controller/consul"
	"github.com/nicholasjackson/smi-controller-sdk/sdk"
	"github.com/nicholasjackson/smi-controller-sdk/sdk/controller"
)

func main() {
	// register our lifecycle callbacks with the controller
	sdk.API().RegisterV1Alpha(&consul.API{})

	// create and start a the controller
	controller.Start()
}
