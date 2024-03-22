package ui

import (
	"log"

	"github.com/jroimartin/gocui"
)

type view struct {
	name           string
	title          string
	x0, y0, x1, y1 int
}

func createActivityView(gui *gocui.Gui) error {
	views := []view{
		{
			title: "addNewActivity",
			name:  "add new time entry",
			x0:    2,
			y0:    1,
			x1:    163,
			y1:    5,
		},
		{
			title: "setCustomer",
			name:  "customer",
			x0:    3,
			y0:    2,
			x1:    42,
			y1:    4,
		},
		{
			title: "setProject",
			name:  "project",
			x0:    43,
			y0:    2,
			x1:    82,
			y1:    4,
		},
		{
			title: "setService",
			name:  "service",
			x0:    83,
			y0:    2,
			x1:    122,
			y1:    4,
		},
		{
			title: "setDescription",
			name:  "description",
			x0:    123,
			y0:    2,
			x1:    162,
			y1:    4,
		},
	}

	for _, view := range views {
		v, err := gui.SetView(view.title, view.x0, view.y0, view.x1, view.y1)
		if err != nil && err != gocui.ErrUnknownView {
			return err
		}

		v.Title = view.title
		v.Editable = true
	}

	if err := gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, activityNextInput); err != nil {
		log.Fatal(err)
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone, activityPrevInput); err != nil {
		log.Fatal(err)
	}

	if _, err := gui.SetCurrentView(views[1].title); err != nil {
		log.Fatal(err)
	}

	return nil
}

func activityNextInput(gui *gocui.Gui, view *gocui.View) error {
	var err error

	switch view.Name() {
	case "setCustomer":
		_, err = gui.SetCurrentView("setProject")
	case "setProject":
		_, err = gui.SetCurrentView("setService")
	case "setService":
		_, err = gui.SetCurrentView("setDescription")
	case "setDescription":
		_, err = gui.SetCurrentView("setCustomer")
	}

	return err
}

func activityPrevInput(gui *gocui.Gui, view *gocui.View) error {
	return nil
}
