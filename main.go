package main

import (
	"flag"
	"fmt"
	"image/color"
	"math/rand"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/hplot/vgshiny"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

var (
	addrFlag = flag.String("addr", ":5555", "server address:port")
	N        = 100
	xmax     = 400
	ymax     = 400
)

type CellState uint8

const (
	Dead CellState = iota
	Alive
)

type Cell struct {
	i, j  int
	state CellState
}

func (c *Cell) switchState() {
	switch c.state {
	case Alive:
		c.state = Dead
	case Dead:
		c.state = Alive
	}
}

type Grid struct {
	C [][]Cell
}

func NewGrid() *Grid {
	G := &Grid{}
	G.C = make([][]Cell, N)
	for i, _ := range G.C {
		G.C[i] = make([]Cell, N)
		for j := range G.C[i] {
			G.C[i][j].state = Dead
			G.C[i][j].i = i
			G.C[i][j].j = j
		}
	}
	return G
}

func (g *Grid) InitRandom() {
	for i := range g.C {
		for j := range g.C[i] {
			c := &g.C[i][j]
			rand := rand.Float64()
			if rand < 0.5 {
				c.state = Dead
			} else {
				c.state = Alive
			}
		}
	}
}

// The video referred to is
//   https://www.youtube.com/watch?v=S-W0NX97DB0
func (g *Grid) InitFirstExampleVideo() {
	g.C[N/2][N/2].state = Alive
	g.C[N/2+1][N/2].state = Alive
	g.C[N/2][N/2+1].state = Alive
	g.C[N/2][N/2+2].state = Alive
	g.C[N/2][N/2+3].state = Alive
	g.C[N/2][N/2-1].state = Alive
}

func (g *Grid) InitClignotant() {
	fmt.Println("N/2", N/2, N/2+1, N/2-1)
	g.C[N/2][N/2].state = Alive
	g.C[N/2][N/2+1].state = Alive
	g.C[N/2][N/2-1].state = Alive
}

func (g *Grid) InitRuche() {
	g.C[N/2][N/2].state = Alive
	g.C[N/2][N/2+1].state = Alive
	g.C[N/2][N/2-1].state = Alive
	g.C[N/2][N/2-2].state = Alive
}

func (g *Grid) Init4Clignotants() {
	g.C[N/2][N/2].state = Alive
	g.C[N/2][N/2+2].state = Alive
	g.C[N/2][N/2+1].state = Alive
	g.C[N/2][N/2-1].state = Alive
	g.C[N/2][N/2-2].state = Alive
}

func (g *Grid) InitDie() {
	g.C[N/2][N/2].state = Alive
	g.C[N/2][N/2+1].state = Alive
	g.C[N/2+1][N/2].state = Alive
	g.C[N/2+1][N/2+1].state = Alive
	g.C[N/2+2][N/2].state = Alive
}

func (g *Grid) Neighbours(c *Cell) []*Cell {
	var neighbours []*Cell
	if c.i > 0 {
		neighbours = append(neighbours, &g.C[c.i-1][c.j])
	}
	if c.i < N-1 {
		neighbours = append(neighbours, &g.C[c.i+1][c.j])
	}
	if c.j > 0 {
		neighbours = append(neighbours, &g.C[c.i][c.j-1])
		if c.i > 0 {
			neighbours = append(neighbours, &g.C[c.i-1][c.j-1])
		}
		if c.i < N-1 {
			neighbours = append(neighbours, &g.C[c.i+1][c.j-1])
		}
	}
	if c.j < N-1 {
		neighbours = append(neighbours, &g.C[c.i][c.j+1])
		if c.i > 0 {
			neighbours = append(neighbours, &g.C[c.i-1][c.j+1])
		}
		if c.i < N-1 {
			neighbours = append(neighbours, &g.C[c.i+1][c.j+1])
		}
	}
	return neighbours
}

func (g *Grid) NoAliveDeadNeighbours(c *Cell) (int, int) { // first: alive, second: dead
	neighbours := g.Neighbours(c)
	var noAlive, noDead int
	for _, n := range neighbours {
		// 		if c.i == 2 {
		// 			fmt.Println("neighbour", n.i, n.j, n.state)
		// 		}
		switch n.state {
		case Alive:
			noAlive++
		case Dead:
			noDead++
		default:
			panic("error !")
		}
	}
	return noAlive, noDead

}

func (g *Grid) Evolve() {
	var cellsToSwitchState []*Cell
	for i := range g.C {
		for j := range g.C[i] {
			c := &g.C[i][j]
			// 			noAliveNeighbours, noDeadNeighbours := g.NoAliveDeadNeighbours(c)
			// 			fmt.Println("numbers ", i, j, " :", noAliveNeighbours, noDeadNeighbours)
			noAliveNeighbours, _ := g.NoAliveDeadNeighbours(c)
			switch {
			case noAliveNeighbours == 3:
				if c.state == Dead {
					cellsToSwitchState = append(cellsToSwitchState, c)
				}
			case noAliveNeighbours < 2 || noAliveNeighbours > 3:
				if c.state == Alive {
					cellsToSwitchState = append(cellsToSwitchState, c)
				}
			}
		}
	}

	for i := range cellsToSwitchState {
		cellsToSwitchState[i].switchState()
	}
}

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
			//fmt.Println("In NewPoints:", c.i, c.j, c.state)
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

/*
func Plot(grid *Grid) {
	points := NewPoints(grid)
	sca, _ := plotter.NewScatter(points)
	sca.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	sca.GlyphStyle.Radius = vg.Points(3.5)
	sca.GlyphStyle.Shape = draw.BoxGlyph{}

	p, _ := plot.New()
	p.X.Min = -0.5
	p.X.Max = float64(N) + 0.5
	p.X.Label.Text = "j"
	p.Y.Min = -0.5
	p.Y.Max = float64(N) + 0.5
	p.Y.Label.Text = "i"
	p.X.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	p.X.Tick.Label.Font.Size = 0
	p.Y.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	p.Y.Tick.Label.Font.Size = 0
	p.Add(sca, plotter.NewGrid())
	p.Save(8*vg.Inch, 8*vg.Inch, "Grid2D.png")

	datac <- Plots{Plot: renderSVG(p)}
}
*/

func Plot(grid *Grid) *hplot.Plot {
	points := NewPoints(grid)
	sca, _ := plotter.NewScatter(points)
	sca.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	sca.GlyphStyle.Radius = vg.Points(2.5)
	sca.GlyphStyle.Shape = draw.BoxGlyph{}

	p := hplot.New()
	p.X.Min = -0.5
	p.X.Max = float64(N) + 0.5
	// 	p.X.Label.Text = "j"
	p.Y.Min = -0.5
	p.Y.Max = float64(N) + 0.5
	// 	p.Y.Label.Text = "i"
	p.X.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	p.X.Tick.Label.Font.Size = 0
	p.X.Tick.Length = 0
	p.Y.Tick.Length = 0
	p.Y.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	p.Y.Tick.Label.Font.Size = 0
	p.Add(sca, plotter.NewGrid())

	return p
}

/*
func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	go webServer(addrFlag)

	///////////////////////////////////////////////////////////////
	// Simple example of grid construction and initialization
	grid := NewGrid()
	//grid.InitToto()
	grid.InitFirstExampleVideo()
	for i := 0; ; i++ {
		fmt.Println("step", i)
		time.Sleep(100 * time.Millisecond)
		grid.Evolve()
		Plot(grid)
	}
	///////////////////////////////////////////////////////////////
}
*/

func display(scr screen.Screen) {

}

func main() {
	flag.Parse()
	grid := NewGrid()
	// 	grid.InitFirstExampleVideo()
	// 	grid.InitClignotant()
	// 	grid.InitRuche()
	// 	grid.Init4Clignotants()
	// 	grid.InitDie()
	grid.InitRandom()
	driver.Main(func(scr screen.Screen) {
		// 		w, err := newWidget(scr, image.Point{xmax, ymax})
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		defer w.Release()

		c, err := vgshiny.New(scr, vg.Length(xmax), vg.Length(ymax))
		if err != nil {
			panic(err)
		}

		c.Run(func(e interface{}) bool {
			switch e := e.(type) {
			case key.Event:
				repaint := false
				switch e.Code {
				case key.CodeQ:
					if e.Direction == key.DirPress {
						return false
					}
				case key.CodeSpacebar:
					if e.Direction == key.DirPress {
						p := Plot(grid)
						p.Draw(draw.New(c))
						grid.Evolve()
						repaint = true
					}
				}
				if repaint {
					c.Send(paint.Event{})
				}

			case paint.Event:
				c.Paint()
			}
			return true
		})
	})
}

// type widget struct {
// 	s      screen.Screen
// 	canvas *vgshiny.Canvas
// }

// func newWidget(s screen.Screen, size image.Point) (*widget, error) {
// 	c, err := vgshiny.New(s, vg.Length(size.X), vg.Length(size.Y))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &widget{s: s, canvas: c}, err
// }

// func (w *widget) Release() {
// 	if w.canvas != nil {
// 		w.canvas.Release()
// 		w.canvas = nil
// 	}
// 	w.s = nil
// }
