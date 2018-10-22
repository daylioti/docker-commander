package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"net"
	"strings"
)

type Docker struct {
	client         *client.Client
	context        context.Context
	containers     []types.Container
	Terminals      []string
	terminalsMap   map[string]int
	updateTerminal fn
}

type fn func(*string)

func (d *Docker) Init(ui fn) {
	var err error
	d.updateTerminal = ui
	d.context = context.Background()
	defer d.context.Done()
	// @Todo set version to settings.yml
	d.client, err = client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}
	d.containers, err = d.client.ContainerList(d.context, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	d.terminalsMap = make(map[string]int)
}

func (d *Docker) Exec(cmd string, path []int, container string) *string {
	var uuid string
	for _, i := range path {
		uuid += string(i)
	}
	uuid += "1" // Just avoid empty string.
	if d.Terminals == nil || d.terminalsMap[uuid] == 0 {
		d.Terminals = append(d.Terminals, "Execute -> "+cmd+"\n")
		d.terminalsMap[uuid] = len(d.Terminals)
	}
	d.updateTerminal(&d.Terminals[d.terminalsMap[uuid]-1])
	d.DockerCommandRun(cmd, container, &d.Terminals[d.terminalsMap[uuid]-1])
	return &d.Terminals[d.terminalsMap[uuid]-1]
}

func (d *Docker) DockerCommandRun(command string, container string, terminal *string) {
	var ExecConfig types.ExecConfig
	var Response types.IDResponse
	var err error
	var ContainerID string
	var HResponse types.HijackedResponse
	var c net.Conn
	ExecConfig.Cmd = strings.Split(command, " ")
	ExecConfig.AttachStdout = true
	ContainerID = d.GetContainerId(container)
	Response, err = d.client.ContainerExecCreate(d.context, ContainerID, ExecConfig)
	if err != nil {
		*terminal += "Docker container from image " + container + " not running \n"
		d.updateTerminal(terminal)
		return
	}
	HResponse, err = d.client.ContainerExecAttach(d.context, Response.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		panic(err)
	}
	c = HResponse.Conn

	var terminalSplit []string
	go func() {
		for {
			terminalSplit = strings.Split(*terminal, "\n")
			if len(terminalSplit) > 58 {
				*terminal = strings.Join(terminalSplit[1:], "\n")
			}
			buf := make([]byte, 1024)
			_, err = c.Read(buf)
			if err != nil {
				_ = c.Close()
				*terminal += "Finished ->" + command + "\n"
				*terminal += "Close connection err: " + err.Error() + "\n"
				d.updateTerminal(terminal)
				return
			} else {
				*terminal += string(buf)
			}

			d.updateTerminal(terminal)
		}
	}()
	_ = d.client.ContainerExecStart(d.context, Response.ID, types.ExecStartCheck{Tty: true})
}

func (d *Docker) GetContainerId(name string) string {
	var container string
	for _, c := range d.containers {
		if c.Image == name {
			container = c.ID
		}
	}
	return container
}
