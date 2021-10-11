package consul

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	"github.com/nicholasjackson/consul-smi-controller/consul/client"
	splitv1alpha4 "github.com/servicemeshinterface/smi-controller-sdk/apis/split/v1alpha4"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func setup(t *testing.T) (*API, *client.MockImpl, logr.Logger) {
	mc := &client.MockImpl{}
	mc.On("WriteServiceSplitter", mock.Anything).Return(nil)
	mc.On("DeleteServiceSplitter", mock.Anything).Return(nil)
	mc.On("WriteServiceRoute", mock.Anything).Return(nil)
	mc.On("DeleteServiceRoute", mock.Anything).Return(nil)

	l := logf.Log.WithName("test")

	return &API{mc}, mc, l
}

func TestUpsertTrafficSplitCallsConsul(t *testing.T) {
	a, mc, l := setup(t)

	splitConfig := splitWithTwoBackend()
	_, err := a.UpsertTrafficSplit(context.Background(), nil, l, splitConfig)
	assert.NoError(t, err)

	mc.AssertCalled(t, "WriteServiceSplitter", mock.Anything)
	args := mc.Mock.Calls[0].Arguments

	ss := args.Get(0).(*api.ServiceSplitterConfigEntry)
	assert.Equal(t, "service-splitter", ss.Kind)
	assert.Equal(t, splitConfig.Spec.Service, ss.Name)

	assert.Equal(t, float32(100), ss.Splits[0].Weight)
	assert.Equal(t, float32(900), ss.Splits[1].Weight)
}

func TestDeleteTrafficSplitCallsConsul(t *testing.T) {
	a, mc, l := setup(t)

	splitConfig := splitWithTwoBackend()
	_, err := a.DeleteTrafficSplit(context.Background(), nil, l, splitConfig)
	assert.NoError(t, err)

	mc.AssertCalled(t, "DeleteServiceSplitter", splitConfig.Spec.Service)
}

func splitWithTwoBackend() *splitv1alpha4.TrafficSplit {
	splitWithTwo := splitv1alpha4.TrafficSplit{
		Spec: splitv1alpha4.TrafficSplitSpec{
			Service: "testservice",
			Backends: []splitv1alpha4.TrafficSplitBackend{
				splitv1alpha4.TrafficSplitBackend{
					Service: "v1",
				},
				splitv1alpha4.TrafficSplitBackend{
					Service: "v2",
				},
			},
		},
	}

	splitWithTwo.Spec.Backends[0].Weight = 100
	splitWithTwo.Spec.Backends[1].Weight = 900

	return &splitWithTwo
}
