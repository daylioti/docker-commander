package ui

import (
	"docker-commander/docker"
	"github.com/gizak/termui"
)

type UI struct {
	cmd      *Commands
}

func (ui *UI) Init(configPath string, dockerClient *docker.Docker) {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	termui.Body.AddRows(
		termui.NewRow(),
		termui.NewRow(),
		termui.NewRow(),
	)


	ui.cmd = &Commands{}
	ui.cmd.Init(configPath, dockerClient)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		ui.cmd.Handle("/sys/kbd/q")
	})
	termui.Handle("/sys/kbd/<left>", func(termui.Event) {
		ui.cmd.Handle("/sys/kbd/<left>")
	})
	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		ui.cmd.Handle("/sys/kbd/<up>")
	})
	termui.Handle("/sys/kbd/<right>", func(termui.Event) {
		ui.cmd.Handle("/sys/kbd/<right>")
	})
	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		ui.cmd.Handle("/sys/kbd/<down>")
	})
	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		ui.cmd.Handle("/sys/kbd/<enter>")
	})

	termui.Loop()
}




func StringColor(text string, color string) string {
	return "[" + text + "](" + color + ")"
}
