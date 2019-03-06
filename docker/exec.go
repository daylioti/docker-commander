package docker

import (
	"github.com/docker/docker/api/types"
	"net"
	"strings"
)

type Exec struct {
	dockerClient    *Docker
	//RunningTerminal *[]string
	Terminals       []*TerminalRun
	//terminalsMap    map[string]int
	//terminalHeight  int
	updateTerminal  fn
}


type TerminalRun struct {
	Command string
	Running bool
	ContainerName string
	ContainerId string
	FromImage string
	Output []string
	Id string
}


type fn func()

func (e *Exec) Init(dockerClient *Docker) {
	e.dockerClient = dockerClient
	//e.terminalsMap = make(map[string]int)
}

func (e *Exec) SetTerminalUpdateFn(ui fn) {
	e.updateTerminal = ui
}
//
//func (e *Exec) SetTerminalHeight(height int) {
//	e.terminalHeight = height
//}


func (e *Exec) CommandExecute(term *TerminalRun) {

	//go func() {

	//}
}

//func (e *Exec) CommandExecute(cmd string, path []int, container string) *[]string {
//	var terminal *[]string
//	//terminal = e.getTerminal(path)
//	*terminal = append(*terminal, "Execute -> "+cmd)
//	//e.ChangeTerminal(path)
//	e.commandRun(cmd, container, terminal)
//	return terminal
//}

func (e *Exec) ChangeTerminal(id string) {
	//e.RunningTerminal = e.getTerminal(path)
	e.updateTerminal()
}
//
//func (e *Exec) getTerminal(path []int) *[]string {
//	var uuid string
//	for _, i := range path {
//		uuid += string(i)
//	}
//	uuid += "1" // Just avoid empty string.
//	//if e.Terminals == nil || e.terminalsMap[uuid] == 0 {
//	//	e.Terminals = append(e.Terminals, []string{})
//	//	e.terminalsMap[uuid] = len(e.Terminals)
//	//}
//	//return &e.Terminals[e.terminalsMap[uuid]-1]
//}

func (e *Exec) commandRun(command string, ContainerID string, terminal *[]string) {
	var ExecConfig types.ExecConfig
	var Response types.IDResponse
	var err error
	var HResponse types.HijackedResponse
	var c net.Conn
	ExecConfig.Cmd = strings.Split(command, " ")
	ExecConfig.AttachStdout = true
	Response, err = e.dockerClient.client.ContainerExecCreate(e.dockerClient.context, ContainerID, ExecConfig)
	if err != nil {
		*terminal = append(*terminal, "Docker container from image " + ContainerID + " not running")
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
			e.updateTerminal()
		}
	}()
	_ = e.dockerClient.client.ContainerExecStart(e.dockerClient.context, Response.ID, types.ExecStartCheck{Tty: true})
}


// Get container id by ContainerName or FromImage or ContainerId params.
func (e *Exec) GetContainerId(term *TerminalRun) string {
	containers, err := e.dockerClient.client.ContainerList(e.dockerClient.context, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {

		if term.FromImage != "" && strings.Contains(c.Image, term.FromImage) {
			return c.ID
		} else if term.ContainerId != "" && c.ID == term.ContainerId {
			return c.ID
		} else if term.ContainerName != "" {
			for _, name := range c.Names {
				if strings.Contains(name, term.ContainerName) {
					return c.ID
				}
			}
		}
	}
	return ""
}
