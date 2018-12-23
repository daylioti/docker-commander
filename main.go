package main

import (
	"docker-commander/docker"
	"docker-commander/ui"
	"flag"
	"fmt"
	"github.com/docker/docker/client"
	"os"
)

var (
	version = "1.0.2"
)

func main() {

	var (
		clientWithVersion = flag.String("api-v", "", "docker api version")
		clientWithHost    = flag.String("api-host", "", "docker api host")
		versionFlag     = flag.Bool("v", false, "output version information and exit")
		helpFlag        = flag.Bool("h", false, "display this help dialog")
	    configFileFlag  = flag.String("c", "", "path to yml config file, default - ./config.yml")
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

	if *clientWithVersion != "" {
      	ops = append(ops, client.WithVersion(*clientWithVersion))
	}

    if *clientWithHost != "" {
    	ops = append(ops, client.WithHost(*clientWithHost))
	}
	dockerClient.Init(ops...)

	UI := new(ui.UI)
	UI.Init(*configFileFlag, dockerClient)
}


var help = `docker-commander - execute commands in docker containers
usage: docker-commander [options]
options:
`

func printHelp() {
	fmt.Println(help)
	flag.PrintDefaults()
}
