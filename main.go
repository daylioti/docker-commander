package main

import (
	"docker-commander/ui"
	"flag"
	"fmt"
	"os"
)

var (
	version = "1.0.2"
)

func main() {

	var (
		versionFlag     = flag.Bool("v", false, "output version information and exit")
		helpFlag        = flag.Bool("h", false, "display this help dialog")
	    configFileFlag  = flag.String("c", "", "path to yml config file, default - ./config.yml")
	)
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

	UI := new(ui.UI)
	UI.Init(*configFileFlag)
}


var help = `docker-commander - execute commands in docker containers
usage: docker-commander [options]
options:
`

func printHelp() {
	fmt.Println(help)
	flag.PrintDefaults()
}
