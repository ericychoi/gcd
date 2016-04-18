package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var dirList []string
var fileToWrite string

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy+1 < len(dirList) {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func changeDir(g *gocui.Gui, v *gocui.View) error {
	var err error
	var l string

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return gocui.ErrQuit
	}
	b := []byte(l)

	if err = ioutil.WriteFile(fileToWrite, b, 0666); err != nil {
		log.Panicln(err)
	}

	return gocui.ErrQuit
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, changeDir); err != nil {
		return err
	}
	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("main", -1, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true

		for i := range dirList {
			fmt.Fprintf(v, "%s\n", dirList[i])
		}

		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fileToWrite = os.Args[1]

	dirs, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for i := range dirs {
		if dirs[i].IsDir() {
			dirList = append(dirList, dirs[i].Name())
		}
	}

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.Cursor = true

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
