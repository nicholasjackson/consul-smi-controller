package client

import (
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/mock"
)

type MockImpl struct {
	mock.Mock
}

func (c *MockImpl) WriteServiceSplitter(ss *api.ServiceSplitterConfigEntry) error {
	args := c.Called(ss)

	return args.Error(0)
}

func (c *MockImpl) DeleteServiceSplitter(name string) error {
	args := c.Called(name)

	return args.Error(0)
}
