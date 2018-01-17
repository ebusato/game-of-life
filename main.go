// TO DO
//  - add text showing iteration

package main

import (
	"flag"
	"fmt"
	"math/rand"

	"golang.org/x/exp/shiny/driver"
)

var (
	addrFlag = flag.String("addr", ":5555", "server address:port")
	N        = 100
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

var grid *Grid

func init() {
	grid = NewGrid()
	grid.InitRandom()
}

func main() {
	flag.Parse()

	driver.Main(GridGraph)
}
