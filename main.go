package main

import (
	"flag"
	"math/rand"
	"time"
)

var (
	addrFlag = flag.String("addr", ":5555", "server address:port")
	N        = 100
)

type CellState uint8

const (
	Alive CellState = iota
	Dead
)

type Cell struct {
	i, j  int
	state CellState
}

func (c *Cell) IsAlive() bool {
	switch c.state {
	case Alive:
		return true
	case Dead:
		return false
	default:
		panic("impossible !")
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
	}
	return G
}

func (g *Grid) InitRandom() {
	for i := range g.C {
		for j := range g.C[i] {
			c := g.C[i][j]
			c.i = i
			c.j = j
			rand := rand.Float64()
			if rand < 0.5 {
				c.state = Dead
			} else {
				c.state = Alive
			}
		}
	}
}

func (g *Grid) PickRandomCell() *Cell {
	i := rand.Intn(N)
	j := rand.Intn(N)
	return &g.C[i][j]
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
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	go webServer(addrFlag)

	///////////////////////////////////////////////////////////////
	// Simple example of grid construction and initialization
	grid := NewGrid(N, 1, 1)
	grid.Init()
	Plot(grid, nil, nil, nil, nil)
	///////////////////////////////////////////////////////////////
}
