// This file contains helper functions to perform web-based plotting

package main

import (
	"image/color"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/vgshiny"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Points struct {
	N int
	X []float64
	Y []float64
}

func NewPoints(g *Grid) *Points {
	points := &Points{}
	for i := range g.C {
		for j := range g.C[i] {
			c := &g.C[i][j]
			if c.state == Alive {
				points.N += 1
				points.X = append(points.X, float64(c.j)) // column
				points.Y = append(points.Y, float64(c.i)) // row
			}
		}
	}
	return points
}

func (p *Points) Len() int {
	return p.N
}

func (p *Points) XY(i int) (x, y float64) {
	return p.X[i], p.Y[i]
}

func setAxisStyle(a *plot.Axis) {
	a.Min = -0.5
	a.Max = float64(N) + 0.5
	a.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	a.Tick.Label.Font.Size = 0
	a.Tick.Length = 0
}

func Plot(g *Grid) *hplot.Plot {
	points := NewPoints(g)
	sca, _ := plotter.NewScatter(points)
	sca.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	sca.GlyphStyle.Radius = vg.Points(3)
	sca.GlyphStyle.Shape = draw.BoxGlyph{}

	p := hplot.New()
	setAxisStyle(&p.X)
	setAxisStyle(&p.Y)
	p.Add(sca, plotter.NewGrid())

	return p
}

func GridGraph(scr screen.Screen) {
	c, err := vgshiny.New(scr, 700, 700)
	if err != nil {
		panic(err)
	}

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
					p := Plot(grid)
					p.Draw(draw.New(c))
					c.Send(paint.Event{})
					grid.Evolve()
				}
			}
		case paint.Event:
			c.Paint()
		}
		return true
	})
}
