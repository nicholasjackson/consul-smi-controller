package consul

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	specsv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha4"
	ctrl "sigs.k8s.io/controller-runtime"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (a *API) UpsertHTTPRouteGroup(
	ctx context.Context,
	r controllerclient.Client,
	l logr.Logger,
	rg *specsv1alpha4.HTTPRouteGroup) (ctrl.Result, error) {

	l.Info("Upsert new HTTPRouteGroup",
		"service", rg.ObjectMeta.Name,
		"matches", rg.Spec.Matches,
	)

	sr := &api.ServiceRouterConfigEntry{Kind: api.ServiceRouter}
	sr.Name = rg.ObjectMeta.Name

	sr.Routes = []api.ServiceRoute{}

	for _, m := range rg.Spec.Matches {
		route := api.ServiceRoute{}
		route.Destination = &api.ServiceRouteDestination{Service: rg.ObjectMeta.Name}

		route.Match = &api.ServiceRouteMatch{}
		route.Match.HTTP = &api.ServiceRouteHTTPMatch{}

		for _, h := range m.Methods {
			switch h {
			case http.MethodConnect:
				fallthrough
			case http.MethodDelete:
				fallthrough
			case http.MethodGet:
				fallthrough
			case http.MethodHead:
				fallthrough
			case http.MethodOptions:
				fallthrough
			case http.MethodPost:
				fallthrough
			case http.MethodPut:
				fallthrough
			case http.MethodTrace:
				route.Match.HTTP.Methods = append(route.Match.HTTP.Methods, h)
			case "*":
				// do nothing as default is allow all
			default:
				return ctrl.Result{}, fmt.Errorf("HTTP method %s, not supported", h)
			}
		}

		route.Match.HTTP.PathRegex = m.PathRegex

		route.Match.HTTP.Header = []api.ServiceRouteHTTPMatchHeader{}
		for k, v := range m.Headers {
			route.Match.HTTP.Header = append(route.Match.HTTP.Header, api.ServiceRouteHTTPMatchHeader{Name: k, Exact: v})
		}

		sr.Routes = append(sr.Routes, route)
	}

	err := a.client.WriteServiceRoute(sr)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// DeleteHTTPRouteGroup implements the API interface method
// for callbacks when a TrafficSplit resource is deleted
// in the Kubernetes cluster
func (a *API) DeleteHTTPRouteGroup(
	ctx context.Context,
	r controllerclient.Client,
	l logr.Logger,
	tt *specsv1alpha4.HTTPRouteGroup) (ctrl.Result, error) {

	l.Info("Delete new HTTPRouteGroup",
		"service", tt.ObjectMeta.Name,
	)

	err := a.client.DeleteServiceRoute(tt.ObjectMeta.Name)
	if err != nil {
		l.Error(err, "Unable to delete ServiceRoute from Consul")
	}

	return ctrl.Result{}, nil
}
