package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v10"
	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	"github.com/nicholasjackson/consul-smi-controller/consul"
	"github.com/nicholasjackson/smi-controller-sdk/sdk"
	"github.com/nicholasjackson/smi-controller-sdk/sdk/controller"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	splitv1alpha1 "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/split/v1alpha1"
	splitClientSet "github.com/servicemeshinterface/smi-sdk-go/pkg/gen/client/split/clientset/versioned"
)

var opts = &godog.Options{
	Format: "pretty",
	Output: colors.Colored(os.Stdout),
}

var logger logr.Logger
var consulClient *api.Client

// store a reference to any objects submitted to the controller for later cleanup
var trafficSplits []*splitv1alpha1.TrafficSplit
var serviceResolvers []*api.ServiceResolverConfigEntry
var serviceDefaults []*api.ServiceConfigEntry

func main() {
	godog.BindFlags("godog.", flag.CommandLine, opts)
	flag.Parse()

	status := godog.TestSuite{
		Name:                "SDK Functional Tests",
		ScenarioInitializer: initializeSuite,
		Options:             opts,
	}.Run()

	os.Exit(status)
}

func initializeSuite(ctx *godog.ScenarioContext) {
	trafficSplits = []*splitv1alpha1.TrafficSplit{}
	serviceResolvers = []*api.ServiceResolverConfigEntry{}
	serviceDefaults = []*api.ServiceConfigEntry{}

	logger = Log()

	config := api.DefaultConfig()
	config.Address = os.Getenv("CONSUL_HTTP_ADDR")

	var err error
	consulClient, err = api.NewClient(config)
	if err != nil {
		panic(err)
	}

	ctx.Step(`^the server is running$`, theServerIsRunning)
	ctx.Step(`^I create a TrafficSplitter$`, iCreateATrafficSplitter)
	ctx.Step(`^the following Consul ServiceResolver exists$`, theFollowingConsulServiceResolverExists)
	ctx.Step(`^I expect a ServiceSplitter to have been created in Consul$`, iExpectAServiceSplitterToHaveBeenCreatedInConsul)
	ctx.Step(`^the following Consul ServiceDefaults exists$`, theFollowingConsulServiceDefaultsExists)

	ctx.AfterScenario(func(s *messages.Pickle, err error) {
		cleanupTrafficSplit()
		cleanupServiceResolvers()
		cleanupServiceDefaults()

		if err != nil {
			fmt.Println(logger.(*StringLogger).String())
		}
	})
}

func cleanupServiceResolvers() {
	for _, ts := range serviceResolvers {
		_, err := consulClient.ConfigEntries().Delete("service-resolver", ts.Name, nil)
		if err != nil {
			fmt.Println(fmt.Errorf("Unable to delete service resolver: %s", err))
		}
	}
}

func cleanupServiceDefaults() {
	for _, ts := range serviceDefaults {
		_, err := consulClient.ConfigEntries().Delete("service-defaults", ts.Name, nil)
		if err != nil {
			fmt.Println(fmt.Errorf("Unable to delete service defaults: %s", err))
		}
	}
}

func cleanupTrafficSplit() {
	c := getK8sConfig()
	sl, err := splitClientSet.NewForConfig(c)
	if err != nil {
		panic(err.Error())
	}

	for _, ts := range trafficSplits {
		sl.SplitV1alpha1().TrafficSplits("default").Delete(context.Background(), ts.Name, v1.DeleteOptions{})

		// also cleanup in consul as we are killing the server before this resolves
		_, err := consulClient.ConfigEntries().Delete("service-splitter", ts.Name, nil)
		if err != nil {
			fmt.Println(fmt.Errorf("Unable to delete service splitter: %s", err))
		}
	}
}

func theServerIsRunning() error {
	api, err := consul.New(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		panic(err)
	}
	sdk.API().RegisterV1Alpha(api)

	// create and start the controller
	config := controller.DefaultConfig()
	config.WebhooksEnabled = false
	config.Logger = logger

	go controller.Start(config)

	return waitForComplete(
		30*time.Second,
		func() error {
			resp, err := http.Get(fmt.Sprintf("http://%s/readyz", config.HealthProbeBindAddress))
			if err == nil {
				if resp != nil && resp.StatusCode == http.StatusOK {
					return nil
				}
			}

			return fmt.Errorf("Timeout waiting for service to become ready")
		},
	)
}

func iCreateATrafficSplitter() error {
	c := getK8sConfig()
	sl, err := splitClientSet.NewForConfig(c)
	if err != nil {
		return err
	}

	ts := &splitv1alpha1.TrafficSplit{
		ObjectMeta: v1.ObjectMeta{Name: "testing"},
		Spec: splitv1alpha1.TrafficSplitSpec{
			Service: "myService",
			Backends: []splitv1alpha1.TrafficSplitBackend{
				splitv1alpha1.TrafficSplitBackend{
					Service: "api-primary",
					Weight:  resource.NewQuantity(100, resource.BinarySI),
				},
			},
		},
	}

	// add to our collection so we can cleanup later
	trafficSplits = append(trafficSplits, ts)

	ts, err = sl.SplitV1alpha1().TrafficSplits("default").Create(context.Background(), ts, v1.CreateOptions{})

	return err
}

func theFollowingConsulServiceResolverExists(arg1 *messages.PickleStepArgument_PickleDocString) error {

	sr := &api.ServiceResolverConfigEntry{}
	serviceResolvers = append(serviceResolvers, sr)

	err := json.Unmarshal([]byte(arg1.Content), sr)
	if err != nil {
		return err
	}

	_, _, err = consulClient.ConfigEntries().Set(sr, nil)
	if err != nil {
		return err
	}

	return nil
}

func theFollowingConsulServiceDefaultsExists(arg1 *messages.PickleStepArgument_PickleDocString) error {

	sr := &api.ServiceConfigEntry{}
	serviceDefaults = append(serviceDefaults, sr)

	err := json.Unmarshal([]byte(arg1.Content), sr)
	if err != nil {
		return err
	}

	_, _, err = consulClient.ConfigEntries().Set(sr, nil)
	if err != nil {
		return err
	}

	return nil
}

func iExpectAServiceSplitterToHaveBeenCreatedInConsul() error {
	return waitForComplete(30*time.Second, func() error {
		_, _, err := consulClient.ConfigEntries().Get("service-splitter", "myService", nil)
		return err
	})
}

func getK8sConfig() *rest.Config {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err.Error())
	}

	return config
}

// helper function to loop until a condition is met
func waitForComplete(duration time.Duration, f func() error) error {
	// wait for the server to mark it is ready
	done := make(chan struct{})
	timeout := time.After(30 * time.Second)

	var err error

	go func() {
		for {
			err = f()
			if err == nil {
				done <- struct{}{}
			}

			time.Sleep(2 * time.Second)
		}
	}()

	select {
	case <-timeout:
		return err
	case <-done:
		return nil
	}
}
