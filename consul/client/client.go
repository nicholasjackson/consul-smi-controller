package client

import (
	"github.com/hashicorp/consul/api"
)

// Consul defines an interface for the Consul client
type Consul interface {
	WriteServiceSplitter(ss *api.ServiceSplitterConfigEntry) error
	DeleteServiceSplitter(name string) error
	WriteServiceRoute(ss *api.ServiceRouterConfigEntry) error
	DeleteServiceRoute(name string) error
}

type consulImpl struct {
	cc *api.Client
}

// New creates a concrete implementation of the Consul interface
func New(addr string) (Consul, error) {
	config := api.DefaultConfig()
	config.Address = addr

	cc, err := api.NewClient(config)
	if err != nil {
		return nil, nil
	}

	return &consulImpl{cc}, nil
}

func (c *consulImpl) WriteServiceSplitter(ss *api.ServiceSplitterConfigEntry) error {
	_, _, err := c.cc.ConfigEntries().Set(ss, nil)

	return err
}

func (c *consulImpl) DeleteServiceSplitter(name string) error {
	_, err := c.cc.ConfigEntries().Delete(api.ServiceSplitter, name, nil)

	return err
}

func (c *consulImpl) WriteServiceRoute(ss *api.ServiceRouterConfigEntry) error {
	_, _, err := c.cc.ConfigEntries().Set(ss, nil)

	return err
}

func (c *consulImpl) DeleteServiceRoute(name string) error {
	_, err := c.cc.ConfigEntries().Delete(api.ServiceRouter, name, nil)

	return err
}
