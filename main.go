package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/daylioti/docker-commander/config"
	"github.com/daylioti/docker-commander/docker"
	"github.com/daylioti/docker-commander/ui"
	"github.com/daylioti/docker-commander/version"
	"github.com/docker/docker/client"
	"github.com/gizak/termui/v3"
)

func main() {

	var (
		configFileFlag = flag.String("c", "", "system path to yml config file or url, default - ./config.yml")
		clientWithHost = flag.String("api-host", "", "docker api host. Example: tcp://127.0.0.1:2376")
		tty            = flag.Bool("tty", false, "enable docker exec tty option")
		color          = flag.Bool("color", false, "display ANSI colors in command output.")
		versionFlag    = flag.Bool("v", false, "output version information and exit")
		helpFlag       = flag.Bool("h", false, "display this help dialog")
	)
	var ops []client.Opt
	flag.Parse()

	if *versionFlag {
		fmt.Println(version.Version)
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
	if *clientWithHost != "" {
		ops = append(ops, client.WithHost(*clientWithHost))
	}
	dockerClient.Init(ops...)
	dockerClient.Exec.Tty = *tty
	dockerClient.Exec.Color = *color

	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	UI := new(ui.UI)
	UI.Init(Cnf, dockerClient, CnfUi)

	uiEvents := termui.PollEvents()
	searchBox := false
	for e := range uiEvents {
		switch true {
		case e.ID == "<MouseLeft>" || e.ID == "<MouseRight>" || e.ID == "<MouseMiddle>" || e.ID == "<MouseRelease>":
			continue
		case e.ID == "<C-s>" || searchBox && e.ID == "<Enter>":
			searchBox = !searchBox
			if !searchBox {
				UI.Commands.Search.Reset()
			} else {
				UI.Commands.Search.Search()
			}
			UI.Render()
		case searchBox:
			UI.Commands.Search.Handle(e.ID)
			if len(UI.Commands.Search.Input) == 0 {
				searchBox = false
			}
		case e.ID == "q" || e.ID == "<C-c>" || e.ID == "Q":
			if len(UI.Commands.Input.Fields) > 0 && e.ID != "<C-c>" {
				UI.Handle(e.ID)
			} else {
				return
			}

		case e.ID == "<Resize>":
			UI.Init(Cnf, dockerClient, CnfUi)
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
