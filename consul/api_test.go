package consul

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	"github.com/nicholasjackson/consul-smi-controller/consul/client"
	splitv1alpha1 "github.com/servicemeshinterface/smi-sdk-go/pkg/apis/split/v1alpha1"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/resource"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

func setup(t *testing.T) (*API, *client.MockImpl, logr.Logger) {
	mc := &client.MockImpl{}
	mc.On("WriteServiceSplitter", mock.Anything).Return(nil)
	mc.On("DeleteServiceSplitter", mock.Anything).Return(nil)

	l := logf.Log.WithName("test")

	return &API{mc}, mc, l
}

func TestUpsertCallsConsul(t *testing.T) {
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

func TestDeleteCallsConsul(t *testing.T) {
	a, mc, l := setup(t)

	splitConfig := splitWithTwoBackend()
	_, err := a.DeleteTrafficSplit(context.Background(), nil, l, splitConfig)
	assert.NoError(t, err)

	mc.AssertCalled(t, "DeleteServiceSplitter", splitConfig.Spec.Service)
}

func splitWithTwoBackend() *splitv1alpha1.TrafficSplit {
	splitWithTwo := splitv1alpha1.TrafficSplit{
		Spec: splitv1alpha1.TrafficSplitSpec{
			Service: "testservice",
			Backends: []splitv1alpha1.TrafficSplitBackend{
				splitv1alpha1.TrafficSplitBackend{
					Service: "v1",
				},
				splitv1alpha1.TrafficSplitBackend{
					Service: "v2",
				},
			},
		},
	}

	w1, _ := resource.ParseQuantity("100m")
	w2, _ := resource.ParseQuantity("900m")

	splitWithTwo.Spec.Backends[0].Weight = &w1
	splitWithTwo.Spec.Backends[1].Weight = &w2

	return &splitWithTwo
}
