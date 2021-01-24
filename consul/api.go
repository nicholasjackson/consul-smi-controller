package consul

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/nicholasjackson/consul-smi-controller/consul/client"
	splitv1alpha1 "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/split/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type API struct {
	client client.Consul
}

// New creates a new API which implements the various
// SDK interfaces called by the Controller
func New(addr string) (*API, error) {
	cc, err := client.New(addr)
	if err != nil {
		return nil, err
	}

	return &API{cc}, nil
}

// UpsertTrafficSplit implements the API interface method
// for callbacks when a TrafficSplit resource is updated or
// created in the Kubernetes cluster
func (a *API) UpsertTrafficSplit(
	ctx context.Context,
	r controllerclient.Client,
	l logr.Logger,
	tt *splitv1alpha1.TrafficSplit) (ctrl.Result, error) {

	return ctrl.Result{}, nil
}
