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
	version = "1.1.0"
)

func main() {

	var (
		clientWithVersion = flag.String("api-v", "", "docker api version")
		clientWithHost    = flag.String("api-host", "", "docker api host")
		versionFlag       = flag.Bool("v", false, "output version information and exit")
		helpFlag          = flag.Bool("h", false, "display this help dialog")
		configFileFlag    = flag.String("c", "", "system path to yml config file or url, default - ./config.yml")
	)
	var ops []func(*client.Client) error
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
	Cnf.Init(*configFileFlag)
	if *clientWithVersion != "" {
		ops = append(ops, client.WithVersion(*clientWithVersion))
	} else {
		ops = append(ops, client.WithVersion("1.39"))
	}

	if *clientWithHost != "" {
		ops = append(ops, client.WithHost(*clientWithHost))
	}
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	dockerClient.Init(ops...)

	UI := new(ui.UI)
	UI.Init(Cnf, dockerClient)

	uiEvents := termui.PollEvents()
	for e := range uiEvents {
		switch e.ID {
		case "q", "<C-c>", "Q":
			return
		case "<Resize>":
			payload := e.Payload.(termui.Resize)
			UI.Grid.SetRect(0, 0, payload.Width, payload.Height)
			UI.Cmd.UpdateRenderElements(Cnf)
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
