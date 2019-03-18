package docker

import (
	"github.com/daylioti/docker-commander/config"
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/docker/docker/api/types"
	"github.com/gizak/termui/v3/widgets"
	"net"
	"strings"
)

type Exec struct {
	dockerClient   *Docker
	Terminals      []*TerminalRun
	updateTerminal fn
}

type TerminalRun struct {
	TabItem     commanderWidgets.TabItem // tab item text and styles
	List        *widgets.List            // list widget with command output
	Active      bool                     // opened tab or not
	Command     string
	Running     bool
	ContainerID string
	ID          string // based on selected item path in menu
	Name        string
	WorkDir     string
}

type fn func(string, bool)

func (e *Exec) Init(dockerClient *Docker) {
	e.dockerClient = dockerClient
}

func (e *Exec) SetTerminalUpdateFn(ui fn) {
	e.updateTerminal = ui
}

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

	go func() {
		for {
			buf := make([]byte, 512)
			_, err = c.Read(buf)
			if err != nil {
				_ = c.Close()
				e.execReadFinish(term, err.Error())
				return
			}
			e.execReadBuffer(term, string(buf))
			e.updateTerminal(term.ID, false)
		}
	}()
	_ = e.dockerClient.client.ContainerExecStart(e.dockerClient.context, Response.ID, types.ExecStartCheck{Tty: true})
}

func (e *Exec) execReadBuffer(term *TerminalRun, buf string) {
	if buf != "" {
		term.List.Rows = append(term.List.Rows, strings.Split(buf, "\n")...)
	}
}

func (e *Exec) execReadFinish(term *TerminalRun, buf string) {
	e.execReadBuffer(term, "Finished -> "+term.Command)
	e.execReadBuffer(term, buf)
	e.updateTerminal(term.ID, true)
}

func (e *Exec) GetTerminal(id string) *TerminalRun {
	for i := 0; i < len(e.Terminals); i++ {
		if e.Terminals[i].ID == id {
			return e.Terminals[i]
		}
	}
	return nil
}

func (e *Exec) GetActiveTerminalIndex() int {
	for i, term := range e.Terminals {
		if term.Active {
			return i
		}
	}
	return -1
}

// Get container id by ContainerName or FromImage or ContainerID params.
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
