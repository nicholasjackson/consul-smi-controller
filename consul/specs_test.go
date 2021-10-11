package consul

import (
	"context"
	"testing"

	"github.com/hashicorp/consul/api"
	specsv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/specs/v1alpha4"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
)

func TestUpsertHTTPRouteGroup(t *testing.T) {
	a, mc, l := setup(t)

	specConfig := httpRouteGroupWithTwoMatches()
	_, err := a.UpsertHTTPRouteGroup(context.Background(), nil, l, specConfig)
	assert.NoError(t, err)

	mc.AssertCalled(t, "WriteServiceRoute", mock.Anything)
	args := mc.Mock.Calls[0].Arguments

	ss := args.Get(0).(*api.ServiceRouterConfigEntry)

	assert.Len(t, ss.Routes, 2)

	assert.Equal(t, specConfig.Spec.Matches[0].Methods, ss.Routes[0].Match.HTTP.Methods)
	assert.Equal(t, specConfig.Spec.Matches[0].PathRegex, ss.Routes[0].Match.HTTP.PathRegex)

	for _, h := range ss.Routes[0].Match.HTTP.Header {
		assert.Equal(t, specConfig.Spec.Matches[0].Headers[h.Name], h.Exact)
	}

	assert.Equal(t, specConfig.Spec.Matches[1].Methods, ss.Routes[1].Match.HTTP.Methods)
	assert.Equal(t, specConfig.Spec.Matches[1].PathRegex, ss.Routes[1].Match.HTTP.PathRegex)

	assert.Len(t, ss.Routes[1].Match.HTTP.Header, 0)
}

func TestUpsertHTTPRouteGroupWithInvalidMethodReturnsError(t *testing.T) {
	a, _, l := setup(t)

	specConfig := httpRouteGroupWithTwoMatches()
	specConfig.Spec.Matches[0].Methods[0] = "*"

	_, err := a.UpsertHTTPRouteGroup(context.Background(), nil, l, specConfig)
	assert.Error(t, err)
}

func TestDeleteHTTPRouteGroupCallsConsul(t *testing.T) {
	a, mc, l := setup(t)

	splitConfig := splitWithTwoBackend()
	_, err := a.DeleteTrafficSplit(context.Background(), nil, l, splitConfig)
	assert.NoError(t, err)

	mc.AssertCalled(t, "DeleteServiceRoute", splitConfig.Spec.Service)
}

func httpRouteGroupWithTwoMatches() *specsv1alpha4.HTTPRouteGroup {
	rg := &specsv1alpha4.HTTPRouteGroup{
		Spec: specsv1alpha4.HTTPRouteGroupSpec{
			Matches: []specsv1alpha4.HTTPMatch{},
		},
	}

	match1 := specsv1alpha4.HTTPMatch{}
	match1.Headers = specsv1alpha4.HTTPHeaders{
		"foo":  "bar",
		"fizz": "buzz",
	}

	match1.Methods = []string{string(specsv1alpha4.HTTPRouteMethodGet)}
	match1.Name = "route1"
	match1.PathRegex = "/.*"

	match2 := specsv1alpha4.HTTPMatch{}

	match2.Methods = []string{string(specsv1alpha4.HTTPRouteMethodPost)}
	match2.Name = "route2"
	match2.PathRegex = "/route2/abc"

	rg.Spec.Matches = append(rg.Spec.Matches, match1)
	rg.Spec.Matches = append(rg.Spec.Matches, match2)

	return rg
}
