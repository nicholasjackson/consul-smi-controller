package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/cucumber/messages-go/v10"
	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/consul-smi-controller/consul"
	"github.com/servicemeshinterface/smi-controller-sdk/sdk"
	"github.com/servicemeshinterface/smi-controller-sdk/sdk/controller"
	"github.com/shipyard-run/shipyard/pkg/clients"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"

	accessv1alpha1 "github.com/servicemeshinterface/smi-controller-sdk/apis/access/v1alpha1"
	accessv1alpha2 "github.com/servicemeshinterface/smi-controller-sdk/apis/access/v1alpha2"
	accessv1alpha3 "github.com/servicemeshinterface/smi-controller-sdk/apis/access/v1alpha3"

	specsv1alpha1 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha1"
	specsv1alpha2 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha2"
	specsv1alpha3 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha3"
	specsv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha4"

	splitv1alpha1 "github.com/servicemeshinterface/smi-controller-sdk/apis/split/v1alpha1"
	splitv1alpha2 "github.com/servicemeshinterface/smi-controller-sdk/apis/split/v1alpha2"
	splitv1alpha3 "github.com/servicemeshinterface/smi-controller-sdk/apis/split/v1alpha3"
	splitv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/split/v1alpha4"
)

var opts = &godog.Options{
	Format: "pretty",
	Output: colors.Colored(os.Stdout),
}

var logger logr.Logger
var k8sClient clients.Kubernetes
var config controller.Config
var consulClient *api.Client

func main() {
	godog.BindFlags("godog.", flag.CommandLine, opts)
	flag.Parse()

	status := godog.TestSuite{
		Name:                 "SDK Functional Tests",
		ScenarioInitializer:  initializeScenario,
		TestSuiteInitializer: initializeSuite,
		Options:              opts,
	}.Run()

	os.Exit(status)
}

func setupClient() error {

	err := setupSplits()
	if err != nil {
		return err
	}

	err = setupAccess()
	if err != nil {
		return err
	}

	err = setupSpecs()
	if err != nil {
		return err
	}

	k8sClient = clients.NewKubernetes(10*time.Millisecond, hclog.NewNullLogger())
	k8sClient, err = k8sClient.SetConfig(os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	config := api.DefaultConfig()
	config.Address = os.Getenv("CONSUL_HTTP_ADDR")
	consulClient, err = api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return err
}

func initializeSuite(ctx *godog.TestSuiteContext) {
	logger = Log()

	err := setupClient()
	if err != nil {
		panic(err)
	}

	config = controller.DefaultConfig()
	config.WebhooksEnabled = true
	config.Logger = logger

	go controller.Start(config)
}

func initializeScenario(ctx *godog.ScenarioContext) {
	// setup the ConsulAPI
	api, err := consul.New(os.Getenv("CONSUL_HTTP_ADDR"))
	if err != nil {
		panic(err)
	}

	// register our lifecycle callbacks with the controller
	sdk.API().RegisterV1Alpha(api)

	ctx.Step(`^the server is running$`, theServerIsRunning)
	ctx.Step(`^I create the following resource$`, iCreateTheFollowingResource)
	ctx.Step(`^I expect "([^"]*)" to be called (\d+) time$`, iExpectToBeCalled)
	ctx.Step(`^I expect the following "([^"]*)" to be created$`, iExpectToBeCreated)
	ctx.Step(`^the following "([^"]*)" config entry exists$`, theConfigEntryExists)

	ctx.AfterScenario(func(s *messages.Pickle, err error) {
		cleanupResources()

		if err != nil {
			fmt.Println("Error occurred running the tests", err)
			fmt.Println(logger.(*StringLogger).String())
		}

		// wait for server to have cleaned up objects and exit
		// as deleting an object is not immediate.
		// We should probably handle this eventual consistency in the cleanup
		// function however this 5 delay should be fine.
		// If you are raising a PR to fix this "should be fine" be sure to shame me in the comments
		time.Sleep(5 * time.Second)
	})
}

func cleanupResources() {
	c := getK8sConfig()
	kc, err := client.New(c, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	kc.DeleteAllOf(
		ctx,
		&accessv1alpha3.TrafficTarget{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v3 TrafficTargets", err)
	}

	kc.DeleteAllOf(
		ctx,
		&accessv1alpha2.TrafficTarget{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v2 TrafficTargets", err)
	}

	kc.DeleteAllOf(
		ctx,
		&accessv1alpha1.TrafficTarget{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v1 TrafficTargets", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha4.HTTPRouteGroup{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v4 HTTPRouteGroup", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha3.HTTPRouteGroup{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v3 HTTPRouteGroup", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha2.HTTPRouteGroup{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v2 HTTPRouteGroup", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha1.HTTPRouteGroup{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v1 HTTPRouteGroup", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha4.TCPRoute{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v4 TCPRoute", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha3.TCPRoute{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v3 TCPRoute", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha2.TCPRoute{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v2 TCPRoute", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha1.TCPRoute{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v1 TCPRoute", err)
	}

	kc.DeleteAllOf(
		ctx,
		&specsv1alpha4.UDPRoute{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v1 UDPRoute", err)
	}

	kc.DeleteAllOf(
		ctx,
		&splitv1alpha1.TrafficSplit{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v1 TrafficSplit", err)
	}

	kc.DeleteAllOf(
		ctx,
		&splitv1alpha2.TrafficSplit{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v2 TrafficSplit", err)
	}

	kc.DeleteAllOf(
		ctx,
		&splitv1alpha3.TrafficSplit{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v3 TrafficSplit", err)
	}

	kc.DeleteAllOf(
		ctx,
		&splitv1alpha4.TrafficSplit{}, client.InNamespace("default"))

	if err != nil {
		fmt.Println("Error removing v4 TrafficSplit", err)
	}

	// clean up consul config
	ce, _, _ := consulClient.ConfigEntries().List("service-defaults", nil)
	for _, c := range ce {
		consulClient.ConfigEntries().Delete(c.GetKind(), c.GetName(), nil)
	}

	ce, _, _ = consulClient.ConfigEntries().List("service-splitter", nil)
	for _, c := range ce {
		consulClient.ConfigEntries().Delete(c.GetKind(), c.GetName(), nil)
	}

	ce, _, _ = consulClient.ConfigEntries().List("service-router", nil)
	for _, c := range ce {
		consulClient.ConfigEntries().Delete(c.GetKind(), c.GetName(), nil)
	}
}

func theServerIsRunning() error {
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

func iCreateTheFollowingResource(arg1 *messages.PickleStepArgument_PickleDocString) error {
	// save the document to a temporary file
	f, err := ioutil.TempFile("", "*.yaml")
	if err != nil {
		return err
	}

	// cleanup
	defer os.Remove(f.Name())

	// write the document to the file
	_, err = f.WriteString(arg1.GetContent())
	if err != nil {
		return err
	}

	// import the file to the kubernetes cluster
	err = k8sClient.Apply([]string{f.Name()}, true)
	return err
}

func iExpectToBeCreated(kind string, resource *messages.PickleStepArgument_PickleDocString) error {
	return waitForComplete(
		30*time.Second,
		func() error {
			switch kind {
			case "service-splitter":
				ss := &api.ServiceSplitterConfigEntry{}

				err := json.Unmarshal([]byte(resource.GetContent()), ss)
				if err != nil {
					return fmt.Errorf("error decoding resource: %s , %v", resource.GetContent(), err)
				}

				// fetch the resource from consul
				ce, _, err := consulClient.ConfigEntries().Get("service-splitter", ss.Name, nil)
				if err != nil {
					return fmt.Errorf("error fetching resource from Consul: %v", err)
				}

				css := ce.(*api.ServiceSplitterConfigEntry)

				if css.Name != ss.Name {
					return fmt.Errorf("expected name to be %s, got %s", ss.Name, css.Name)
				}

				if len(css.Splits) != len(ss.Splits) {
					return fmt.Errorf("expected %d splits, got %d", len(ss.Name), len(css.Name))
				}

				for i, s := range css.Splits {
					if s.Service != css.Splits[i].Service {
						return fmt.Errorf("expected service name %s for split, got %s", s.Service, css.Splits[i].Service)
					}
					if s.Weight != css.Splits[i].Weight {
						return fmt.Errorf("expected service weight %f for split, got %f", s.Weight, css.Splits[i].Weight)
					}
				}
			case "service-router":
				ss := &api.ServiceRouterConfigEntry{}

				err := json.Unmarshal([]byte(resource.GetContent()), ss)
				if err != nil {
					return fmt.Errorf("error decoding resource: %s , %v", resource.GetContent(), err)
				}

				// fetch the resource from consul
				ce, _, err := consulClient.ConfigEntries().Get("service-router", ss.Name, nil)
				if err != nil {
					return fmt.Errorf("error fetching resource from Consul: %v", err)
				}

				css := ce.(*api.ServiceRouterConfigEntry)

				if css.Name != ss.Name {
					return fmt.Errorf("expected name to be %s, got %s", ss.Name, css.Name)
				}

			default:
				return fmt.Errorf("unknown resource kind: %s", kind)
			}

			return nil
		})
}

func theConfigEntryExists(kind string, resource *messages.PickleStepArgument_PickleDocString) error {
	switch kind {
	case "service-defaults":
		ss := &api.ServiceConfigEntry{}
		err := json.Unmarshal([]byte(resource.GetContent()), ss)
		if err != nil {
			return fmt.Errorf("error decoding config entry: %s , %v", resource.GetContent(), err)
		}

		_, _, err = consulClient.ConfigEntries().Set(ss, nil)
		if err != nil {
			return fmt.Errorf("error setting config entry in Consul: %v", err)
		}
	}

	return nil
}

// The controller is eventually consistent so we need to check this in a loop
func iExpectToBeCalled(method string, n int) error {
	return godog.ErrPending
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
	timeout := time.After(duration)

	var err error

	go func() {
		for {
			err = f()
			if err == nil {
				done <- struct{}{}
				break
			}

			// retry after 1s
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-timeout:
		return err
	case <-done:
		return nil
	}
}
