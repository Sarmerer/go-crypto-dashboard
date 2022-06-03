package bucket

import (
	"github.com/sarmerer/go-crypto-dashboard/passivbotmgr"
	"github.com/sarmerer/go-crypto-dashboard/passivbotmgr/droplet"
)

// bucket is an api that allows you to communicate with droplets
// and get information about them.

type Bucket interface {
	Start() error
	Stop() error

	GetDroplets() ([]droplet.Droplet, error)
	GetDroplet(id passivbotmgr.DropletID) (droplet.Droplet, error)
	GetDropletByIP(ip passivbotmgr.DropletIP) (droplet.Droplet, error)
}

type bucket struct {
	droplets []droplet.Droplet
}

func NewBucket() Bucket {
	return &bucket{}
}

func (b *bucket) Start() error {
	return nil
}

func (b *bucket) Stop() error {
	return nil
}

func (b *bucket) GetDroplets() ([]droplet.Droplet, error) {
	return b.droplets, nil
}

func (b *bucket) GetDroplet(id passivbotmgr.DropletID) (droplet.Droplet, error) {
	return nil, nil
}

func (b *bucket) GetDropletByIP(ip passivbotmgr.DropletIP) (droplet.Droplet, error) {
	return nil, nil
}
