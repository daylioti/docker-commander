package docker

import (
	"context"
	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"strconv"
)

// Docker main structure.
type Docker struct {
	client  *client.Client
	context context.Context
	Exec    *Exec
	Ping    types.Ping
}

// Init initialize connection to docker.
func (d *Docker) Init(ops ...client.Opt) {
	var err error
	d.context = context.Background()
	defer d.context.Done()
	if d.client, err = client.NewClientWithOpts(ops...); err != nil {
		panic(err)
	}
	_, err = d.client.Ping(d.context)
	if err != nil {
		panic(err)
	}
	min, _ := strconv.ParseFloat(api.MinVersion, 32)
	clientAPIVersion, _ := strconv.ParseFloat(d.Ping.APIVersion, 32)
	max, _ := strconv.ParseFloat(api.DefaultVersion, 32)
	if min <= clientAPIVersion && clientAPIVersion <= max {
		ops = append(ops, client.WithVersion(d.Ping.APIVersion))
	}
	if d.client, err = client.NewClientWithOpts(ops...); err != nil {
		panic(err)
	}
	d.Exec = &Exec{}
	d.Exec.Init(d)
}
