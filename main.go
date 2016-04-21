package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

// Node represents a directory
type Node struct {
	children []*Node
	dir      string
}

func (n *Node) String() string {
	return n.dir
}

var root *Node
var fileToWrite string
var nodeCount int

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if cy+1 < nodeCount {
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

func writeDir(g *gocui.Gui, v *gocui.View) error {
	var err error
	var l string

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		return gocui.ErrQuit
	}
	b := []byte(strings.TrimSpace(l))

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
	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, writeDir); err != nil {
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
		printTree(v, root, 0)
		if err := g.SetCurrentView("main"); err != nil {
			return err
		}
	}
	return nil
}

func addChildren(n *Node) {
	nodeCount++
	dirs, err := ioutil.ReadDir(n.dir)
	if err != nil {
		log.Fatal(err)
	}

	for i := range dirs {
		if dirs[i].IsDir() {
			child := &Node{dir: fmt.Sprintf("%s/%s", n.dir, dirs[i].Name())}
			n.children = append(n.children, child)
			addChildren(child)
		}
	}
}

func printTree(w io.Writer, n *Node, level int) {
	fmt.Fprintf(w, "%s\n", n)
	for i := 0; i < level; i++ {
		fmt.Fprint(w, ` `)
	}
	for c := range n.children {
		printTree(w, n.children[c], level+1)
	}
}

func main() {
	fileToWrite = os.Args[1]
	root = &Node{dir: "."}
	addChildren(root)

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
