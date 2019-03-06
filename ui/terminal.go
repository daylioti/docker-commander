package ui

import "github.com/daylioti/docker-commander/docker"

type TerminalUi struct {
	ui           *UI
	client       *docker.Docker
    //terminals    []*TerminalRun
}



func (t *TerminalUi) Init(ui *UI, client *docker.Docker) {
	t.ui = ui
    t.client = client
}

func GetIdFromPath(path []int) string {
	id := "0"
	for _, i := range path {
		id += string(i)
	}
	return id
}


func (t *TerminalUi) Execute(term *docker.TerminalRun) {
	//t.client.Exec.Terminals = append(term)
	t.client.Exec.CommandExecute(term)
	//var containerId string
	//containerId = t.client.Exec.GetContainerId(term.Container)
	//if containerId != "" {
	//	term.Running = true
	//	t.client.Exec.CommandExecute(term.Command, containerId, &t.terminals[len(t.terminals)-1].Output)
	//}

}
