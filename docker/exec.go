package docker

import (
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/ui/helpers"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/docker/docker/api/types"
	"github.com/gizak/termui/v3"
	"net"
	"strings"
)

const bufferReadSize = 512

// Exec - main struct for execute commands inside docker containers.
type Exec struct {
	Tty            bool
	Color          bool
	dockerClient   *Docker
	Terminals      []*TerminalRun
	updateTerminal fn
}

// TerminalRun - struct for running commands.
type TerminalRun struct {
	TabItem     *commanderWidgets.TabItem      // tab item text and styles
	List        *commanderWidgets.TerminalList // list widget with command output
	Active      bool                           // opened tab or not
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
	ExecConfig.Tty = e.Tty
	Response, err = e.dockerClient.client.ContainerExecCreate(e.dockerClient.context, term.ContainerID, ExecConfig)
	width, height := termui.TerminalDimensions()
	resize := types.ResizeOptions{
		Width:  uint(width),
		Height: uint(height),
	}
	_ = e.dockerClient.client.ContainerResize(e.dockerClient.context, term.ContainerID, resize)

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
	e.execReadBuffer(term, []byte(commanderWidgets.StyleText("Executed:", "fg:green")), false)
	e.execReadBuffer(term, []byte(commanderWidgets.StyleText("Dir -> "+term.WorkDir, "fg:green")), false)
	e.execReadBuffer(term, []byte(commanderWidgets.StyleText("Cmd -> "+term.Command, "fg:green")), false)
	e.execReadBuffer(term, []byte(commanderWidgets.StyleText("ContainerId -> "+term.ContainerID, "fg:green")), false)

	go func() {
		for {
			b := make([]byte, bufferReadSize)
			if _, err = c.Read(b); err != nil {
				_ = c.Close()
				e.execReadFinish(term, err.Error())
				return
			}
			e.execReadBuffer(term, b, e.Color)
			e.updateTerminal(term, false)
		}
	}()
	_ = e.dockerClient.client.ContainerExecStart(e.dockerClient.context, Response.ID, types.ExecStartCheck{Tty: true})
}

// execReadBuffer - paste result rom output buffer to display list.
func (e *Exec) execReadBuffer(term *TerminalRun, buff []byte, color bool) {
	if len(buff) > 0 {
		if color {
			var rows []string
			for _, line := range strings.Split(string(buff), "\n") {
				rows = strings.Split(line, "\b\b")
				if len(rows) >= 1 && strings.Contains(line, "\b\b") {
					// Remove prev line
					term.List.Rows = term.List.Rows[:len(term.List.Rows)-1]
				}
				line = string(helpers.TTYColorsParse([]byte(rows[len(rows)-1])))
				term.List.Rows = append(term.List.Rows, line)
			}
		} else {
			term.List.Rows = append(term.List.Rows, strings.Split(string(buff), "\n")...)
		}
	}
}

// execReadFinish - paste finish info to display list.
func (e *Exec) execReadFinish(term *TerminalRun, buf string) {
	e.execReadBuffer(term, []byte(commanderWidgets.StyleText("Finished -> "+term.Command, "fg:green")), false)
	if buf != "EOF" {
		e.execReadBuffer(term, []byte(buf), e.Color)
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
