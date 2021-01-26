package consul

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	"github.com/nicholasjackson/consul-smi-controller/consul/client"
	splitv1alpha1 "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/split/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// API defines a controller adaptor for Conul Service Mesh
// it contains all the required callbacks invoked by the SMI controller
// SDK.
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

	l.Info("Upsert new Traffic Split")
	l.Info("details",
		"service", tt.Spec.Service,
		"backends", tt.Spec.Backends,
	)

	ss := &api.ServiceSplitterConfigEntry{Kind: api.ServiceSplitter}
	ss.Name = tt.Spec.Service
	ss.Splits = []api.ServiceSplit{}

	for _, b := range tt.Spec.Backends {
		ss.Splits = append(
			ss.Splits,
			api.ServiceSplit{
				Service:       ss.Name,
				ServiceSubset: b.Service,
				Weight:        float32(b.Weight.AsDec().UnscaledBig().Int64()),
			})
	}

	err := a.client.WriteServiceSplitter(ss)

	// if we can not process the request we should not keep retrying
	// returning a default ctrl.Request{} and no error
	// tells the controller we have accepted the config
	// TODO: we need to determine recoverable errors from non
	// recoverable
	if err != nil {
		l.Error(err, "Unable to write service splitter to Consul")
	}

	return ctrl.Result{}, nil
}

// DeleteTrafficSplit implements the API interface method
// for callbacks when a TrafficSplit resource is deleted
// in the Kubernetes cluster
func (a *API) DeleteTrafficSplit(
	ctx context.Context,
	r controllerclient.Client,
	l logr.Logger,
	tt *splitv1alpha1.TrafficSplit) (ctrl.Result, error) {

	l.Info("Delete new Traffic Split")
	l.Info("details",
		"service", tt.Spec.Service,
		"backends", tt.Spec.Backends,
	)

	// if we can not process the request we should not keep retrying
	// returning a default ctrl.Request{} and no error
	// tells the controller we have accepted the config
	// TODO: we need to determine recoverable errors from non
	// recoverable
	err := a.client.DeleteServiceSplitter(tt.Spec.Service)
	if err != nil {
		l.Error(err, "Unable to delete service splitter from Consul")
	}

	return ctrl.Result{}, nil
}
