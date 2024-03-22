package ui

import (
	"log"

	"github.com/joergmis/track"
	"github.com/jroimartin/gocui"
)

type ui struct {
	repository track.Repository
}

func New(repo track.Repository) track.UI {
	return &ui{
		repository: repo,
	}
}

func (u *ui) Run() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	createActivityView(g)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
