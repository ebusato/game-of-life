// This file contains helper functions to perform web-based plotting

package main

import (
	"go-hep.org/x/hep/hplot/vgshiny"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
	"gonum.org/v1/plot/vg/draw"
)

func GridGraph(scr screen.Screen) {
	c, err := vgshiny.New(scr, 700, 700)
	if err != nil {
		panic(err)
	}
	grid := NewGrid()
	grid.Init()

	c.Run(func(e interface{}) bool {
		switch e := e.(type) {
		case key.Event:
			switch e.Code {
			case key.CodeQ:
				if e.Direction == key.DirPress {
					return false
				}
			case key.CodeSpacebar:
				if e.Direction == key.DirPress {
					grid.Evolve()
					p := grid.Draw()
					p.Draw(draw.New(c))
					c.Send(paint.Event{})
				}
			}
		case paint.Event:
			c.Paint()
		}
		return true
	})
}
