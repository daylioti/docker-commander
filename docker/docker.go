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
	RunningTerminal *string
	Terminals      []string
	terminalsMap   map[string]int
	terminalHeight int
	updateTerminal fn
}

type fn func()

func (d *Docker) Init(ui fn) {
	var err error
	d.updateTerminal = ui
	d.context = context.Background()
	defer d.context.Done()
	// @Todo set version to settings.yml
	d.client, err = client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic("Docker not running.")
	}
	d.terminalsMap = make(map[string]int)
}

func (d *Docker) SetTerminalHeight(height int) {
	d.terminalHeight = height
}

func (d *Docker) Exec(cmd string, path []int, container string) *string {
	var terminal *string
    terminal = d.getTerminal(path)
	*terminal = "Execute -> "+cmd+"\n"
	d.ChangeTerminal(path)
	d.DockerCommandRun(cmd, container, terminal)
	return terminal
}

func (d *Docker) ChangeTerminal(path []int) {
	d.RunningTerminal = d.getTerminal(path)
	d.updateTerminal()
}

func (d *Docker) getTerminal(path []int) *string {
	var uuid string
	for _, i := range path {
		uuid += string(i)
	}
	uuid += "1" // Just avoid empty string.
	if d.Terminals == nil || d.terminalsMap[uuid] == 0 {
		d.Terminals = append(d.Terminals, "")
		d.terminalsMap[uuid] = len(d.Terminals)
	}
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
		d.updateTerminal()
		return
	}
	HResponse, err = d.client.ContainerExecAttach(d.context, Response.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		panic(err)
	}
	c = HResponse.Conn

	go func() {
		for {
			buf := make([]byte, 512)
			_, err = c.Read(buf)
			if err != nil {
				_ = c.Close()
				*terminal += "Finished -> " + command + "\n"
				*terminal += "Close connection err: " + err.Error() + "\n"
				d.updateTerminal()
				return
			} else {
				*terminal += string(buf)
			}
			var text []string
			var sliceTrim int
			text = strings.Split(*terminal, "\n")
			if len(text) > d.terminalHeight {
				sliceTrim = len(text) - d.terminalHeight
				text = append(text[:0], text[sliceTrim:]...)
				*terminal = strings.Join(text, "\n")
			}
			d.updateTerminal()
		}
	}()
	_ = d.client.ContainerExecStart(d.context, Response.ID, types.ExecStartCheck{Tty: true})
}

func (d *Docker) GetContainerId(name string) string {
	containers, _ := d.client.ContainerList(d.context, types.ContainerListOptions{})
	for _, c := range containers {
		if strings.Contains(c.Image, name) {
			return c.ID
		}
	}
	return ""
}
