package docker

import (
	"context"
	"github.com/docker/docker/api"
	"github.com/docker/docker/client"
	"strconv"
)

type Docker struct {
	client  *client.Client
	context context.Context
	Exec    *Exec
}

func (d *Docker) Init(version string, ops ...func(*client.Client) error) {
	var err error
	d.context = context.Background()
	defer d.context.Done()
	if version != "" {
		ops = append(ops, client.WithVersion(version))
	}
	d.client, err = client.NewClientWithOpts(ops...)
	if err != nil {
		panic(err)
	}
	if version == "" {
		ping, err := d.client.Ping(d.context)
		if err != nil {
			panic(err)
		}
		min, _ := strconv.ParseFloat(api.MinVersion, 32)
		clientApiVersion, _ := strconv.ParseFloat(ping.APIVersion, 32)
		max, _ := strconv.ParseFloat(api.DefaultVersion, 32)
		if min <= clientApiVersion && clientApiVersion <= max {
			ops = append(ops, client.WithVersion(ping.APIVersion))
		}
	}
	d.client, err = client.NewClientWithOpts(ops...)
	if err != nil {
		panic(err)
	}
	d.Exec = &Exec{}
	d.Exec.Init(d)
}
