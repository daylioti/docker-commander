package docker

import (
	"github.com/docker/docker/api/types"
	"net"
	"strings"
)

type Exec struct {
	dockerClient    *Docker
	RunningTerminal *[]string
	Terminals       [][]string
	terminalsMap    map[string]int
	terminalHeight  int
	updateTerminal  fn
}
type fn func()

func (e *Exec) Init(dockerClient *Docker) {
	e.dockerClient = dockerClient
	e.terminalsMap = make(map[string]int)
}

func (e *Exec) SetTerminalUpdateFn(ui fn) {
	e.updateTerminal = ui
}

func (e *Exec) SetTerminalHeight(height int) {
	e.terminalHeight = height
}

func (e *Exec) CommandExecute(cmd string, path []int, container string) *[]string {
	var terminal *[]string
	terminal = e.getTerminal(path)
	*terminal = append(*terminal, "Execute -> "+cmd)
	e.ChangeTerminal(path)
	e.commandRun(cmd, container, terminal)
	return terminal
}

func (e *Exec) ChangeTerminal(path []int) {
	e.RunningTerminal = e.getTerminal(path)
	e.updateTerminal()
}

func (e *Exec) getTerminal(path []int) *[]string {
	var uuid string
	for _, i := range path {
		uuid += string(i)
	}
	uuid += "1" // Just avoid empty string.
	if e.Terminals == nil || e.terminalsMap[uuid] == 0 {
		e.Terminals = append(e.Terminals, []string{})
		e.terminalsMap[uuid] = len(e.Terminals)
	}
	return &e.Terminals[e.terminalsMap[uuid]-1]
}

func (e *Exec) commandRun(command string, container string, terminal *[]string) {
	var ExecConfig types.ExecConfig
	var Response types.IDResponse
	var err error
	var ContainerID string
	var HResponse types.HijackedResponse
	var c net.Conn
	ExecConfig.Cmd = strings.Split(command, " ")
	ExecConfig.AttachStdout = true
	ContainerID = e.GetContainerId(container)
	Response, err = e.dockerClient.client.ContainerExecCreate(e.dockerClient.context, ContainerID, ExecConfig)
	if err != nil {
		*terminal = append(*terminal, "Docker container from image "+container+" not running")
		e.updateTerminal()
		return
	}
	HResponse, err = e.dockerClient.client.ContainerExecAttach(e.dockerClient.context, Response.ID, types.ExecStartCheck{Tty: true})
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
				*terminal = append(*terminal, "Finished -> "+command)
				e.updateTerminal()
				return
			} else {
				*terminal = append(*terminal, string(buf))
			}
			//var text []string
			//var sliceTrim int
			//text = strings.Split(*terminal, "\n")
			//if len(text) > e.terminalHeight {
			//	sliceTrim = len(text) - e.terminalHeight
			//	//text = append(text[:0], text[sliceTrim:]...)
			//	*terminal = strings.Join(text[sliceTrim+5:], "\n")
			//}
			e.updateTerminal()
		}
	}()
	_ = e.dockerClient.client.ContainerExecStart(e.dockerClient.context, Response.ID, types.ExecStartCheck{Tty: true})
}

func (e *Exec) GetContainerId(name string) string {
	containers, err := e.dockerClient.client.ContainerList(e.dockerClient.context, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		if strings.Contains(c.Image, name) {
			return c.ID
		}
	}
	return ""
}
