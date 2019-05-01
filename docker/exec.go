package docker

import (
	"github.com/daylioti/docker-commander/config"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/docker/docker/api/types"
	"github.com/gizak/termui/v3/widgets"
	"net"
	"strings"
)

// Exec - main struct for execute commands inside docker containers.
type Exec struct {
	dockerClient   *Docker
	Terminals      []*TerminalRun
	updateTerminal fn
}

// TerminalRun - struct for running commands.
type TerminalRun struct {
	TabItem     *commanderWidgets.TabItem // tab item text and styles
	List        *widgets.List             // list widget with command output
	Active      bool                      // opened tab or not
	Command     string
	Running     bool
	ContainerID string
	ID          string // based on selected item path in menu
	Name        string
	WorkDir     string
}

// Function for re-render after receiving some updates from execution command.
type fn func(*TerminalRun, bool)

// Init initialize exec object.
func (e *Exec) Init(dockerClient *Docker) {
	e.dockerClient = dockerClient
}

// SetTerminalUpdateFn set function for updating terminal list from command output.
func (e *Exec) SetTerminalUpdateFn(ui fn) {
	e.updateTerminal = ui
}

// CommandRun - execute process in docker container.
func (e *Exec) CommandRun(term *TerminalRun) {
	var ExecConfig types.ExecConfig
	var Response types.IDResponse
	var err error
	var HResponse types.HijackedResponse
	var c net.Conn
	if term.ContainerID == "" {
		e.execReadFinish(term, "Can't find running container")
	}
	ExecConfig.Cmd = strings.Split(term.Command, " ")
	ExecConfig.AttachStdout = true
	if term.WorkDir != "" {
		ExecConfig.WorkingDir = term.WorkDir
	}
	ExecConfig.Tty = true
	Response, err = e.dockerClient.client.ContainerExecCreate(e.dockerClient.context, term.ContainerID, ExecConfig)
	if err != nil {
		e.execReadFinish(term, err.Error())
		return
	}
	HResponse, err = e.dockerClient.client.ContainerExecAttach(e.dockerClient.context, Response.ID,
		types.ExecStartCheck{Tty: true})
	if err != nil {
		e.execReadFinish(term, err.Error())
		return
	}
	c = HResponse.Conn
	e.execReadBuffer(term, "[Executed:](fg:green)")
	e.execReadBuffer(term, "[Dir -> "+term.WorkDir+"](fg:green)")
	e.execReadBuffer(term, "[Cmd -> "+term.Command+"](fg:green)")
	e.execReadBuffer(term, "[ContainerId -> "+term.ContainerID+"](fg:green)")

	go func() {
		for {
			buf := make([]byte, 512)
			if _, err = c.Read(buf); err != nil {
				_ = c.Close()
				e.execReadFinish(term, err.Error())
				return
			}
			e.execReadBuffer(term, string(buf))
			e.updateTerminal(term, false)
		}
	}()
	_ = e.dockerClient.client.ContainerExecStart(e.dockerClient.context, Response.ID, types.ExecStartCheck{Tty: true})
}

// execReadBuffer - paste result rom output buffer to display list.
func (e *Exec) execReadBuffer(term *TerminalRun, buf string) {
	if buf != "" {
		term.List.Rows = append(term.List.Rows, strings.Split(buf, "\n")...)
	}
}

// execReadFinish - paste finish info to display list.
func (e *Exec) execReadFinish(term *TerminalRun, buf string) {
	e.execReadBuffer(term, "[Finished -> "+term.Command+"](fg:green)")
	if buf != "EOF" {
		e.execReadBuffer(term, buf)
	}
	e.updateTerminal(term, true)
}

// GetActiveTerminalIndex get array index of active terminal.
func (e *Exec) GetActiveTerminalIndex() int {
	for i, term := range e.Terminals {
		if term.Active {
			return i
		}
	}
	return -1
}

// GetContainerID - get container id by ContainerName or FromImage or ContainerID params.
func (e *Exec) GetContainerID(config config.Config) string {
	containers, err := e.dockerClient.client.ContainerList(e.dockerClient.context, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, c := range containers {
		if config.Exec.Connect.FromImage != "" && strings.Contains(c.Image, config.Exec.Connect.FromImage) {
			return c.ID
		} else if config.Exec.Connect.ContainerID != "" && c.ID == config.Exec.Connect.ContainerID {
			return c.ID
		} else if config.Exec.Connect.ContainerName != "" {
			for _, name := range c.Names {
				if strings.Contains(name, config.Exec.Connect.ContainerName) {
					return c.ID
				}
			}
		}
	}
	return ""
}
