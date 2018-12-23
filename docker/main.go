package docker

import (
	"context"
	"github.com/docker/docker/client"
)

type Docker struct {
	client *client.Client
	context context.Context
	Exec *Exec
}

func (d *Docker) Init(ops ...func(*client.Client) error) {
	var err error
	d.context = context.Background()
	defer d.context.Done()
	d.client, err = client.NewClientWithOpts(ops...)
	if err != nil {
		panic("Docker not running.")
	}
	d.Exec = &Exec{}
	d.Exec.Init(d)
}



