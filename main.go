package main

import (
	"flag"
	"fmt"
	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui"
	"github.com/docker/docker/client"
	"github.com/gizak/termui/v3"
	"os"
)

var (
	version = "1.1.2"
)

func main() {

	var (
		clientWithVersion = flag.String("api-v", "", "docker api version")
		clientWithHost    = flag.String("api-host", "", "docker api host. Example: tcp://127.0.0.1:2376")
		tty               = flag.Bool("tty", false, "Enable docker exec tty option with parse colors")
		versionFlag       = flag.Bool("v", false, "output version information and exit")
		helpFlag          = flag.Bool("h", false, "display this help dialog")
		configFileFlag    = flag.String("c", "", "system path to yml config file or url, default - ./config.yml")
	)
	var ops []client.Opt
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	if *configFileFlag == "" {
		*configFileFlag = "./config.yml"
	}
	dockerClient := &docker.Docker{}
	Cnf := &config.Config{}
	CnfUi := &config.UIConfig{}
	config.CnfInit(*configFileFlag, Cnf, CnfUi)
	Cnf.Init()
	if *clientWithVersion != "" {
		ops = append(ops, client.WithVersion(*clientWithVersion))
	}
	if *clientWithHost != "" {
		ops = append(ops, client.WithHost(*clientWithHost))
	}

	dockerClient.Init(*clientWithVersion, ops...)
	dockerClient.Exec.Tty = *tty

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	UI := new(ui.UI)
	UI.Init(Cnf, dockerClient, CnfUi)

	uiEvents := termui.PollEvents()
	for e := range uiEvents {
		switch e.ID {
		case "q", "<C-c>", "Q":
			if len(UI.Input.Fields) > 0 && e.ID != "<C-c>" {
				UI.Handle(e.ID)
			} else {
				return
			}

		case "<Resize>":
			payload := e.Payload.(termui.Resize)
			UI.TermHeight = payload.Height
			UI.TermWidth = payload.Width
			UI.Render()
		default:
			UI.Handle(e.ID)
		}
	}
}

var help = `docker-commander - execute commands in docker containers
usage: docker-commander [options]
options:
`

func printHelp() {
	fmt.Println(help)
	flag.PrintDefaults()
}
