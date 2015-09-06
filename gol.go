// Implementing Conway's Game Of Life
// ----------------------------------
//
// Using a map for storing the current state of the world.
//
// We are printing the successive populations in a format that can be fed
// to gnuplot and creating in this way an animated view of the population.
//
// This is just an exercise for using maps in go! Do not take this
// too serious...
//
// To see the simulation in gnuplot, call the program like this:
// ./gol.exe | gnuplot --persist

package main

import (
	"fmt"
)

// We are storing the cells (alive or dead) in a map. The keys are the Cartesian
// coordinates of the cells and the values are the properties of the cells,
// namely their state and number of alive neighbours.

// A cell has its state, and its number of life neighbours
type Cell struct {
	alive bool
	n int
}

// The coordinates are plain 2-d cartesian coordinates
type Coord struct {
	x int
	y int
}

// The world is a map of Coord and Cell
type World map[Coord]Cell

// Inflate inflates the world with dead cells surrounding
// the live cells
func (world World) Inflate() World {
	var newWorld World
	newWorld = make(World)
	
	for coord, cell := range world {
		newWorld[coord] = cell
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				c := Coord{coord.x + i, coord.y + j}
				if _, found := newWorld[c]; !found {
					newWorld[c] = Cell{false, 0}
				}
			}
		}
	}

	return newWorld
}

// Deflate deflates the world: only the live cells remain
func (world World) Deflate() World {
	var newWorld World
	newWorld = make(World)
	
	for coord, cell := range world {
		if cell.alive {
			newWorld[coord] = cell
		}
	}

	return newWorld
}

// Tick computes the next generation of live cells in the world
func (world World) Tick() World {
	// count live neighbours for each cell
	for coord, cell := range world {
		n := 0
		for i := -1; i < 2; i++ {
			for j := -1; j < 2; j++ {
				c := Coord{coord.x + i, coord.y + j}
				if (i != 0 || j != 0) && world[c].alive {
					n = n+1
				}
			}
		}
		world[coord] = Cell{cell.alive, n}
	}

	var newWorld World
	newWorld = make(World)
	
	// apply the rules of the game to each cell
	for coord, cell := range world {
		if cell.alive {
			if 1 < cell.n && cell.n < 4 {
				newWorld[coord] = Cell{true, 0}
			}
		} else {
			if cell.n == 3 {
				newWorld[coord] = Cell{true, 0}
			}
		}
	}

	return newWorld
}

// gnuplotHeader prints the header for gnuplot
func gnuplotHeader() {
	fmt.Println("unset key; set xrange[-50:50]")
	fmt.Println("set yrange[-50:50]")
}

// gnuplotCell prints the coordinates of a cell
func gnuplotWorld(world World) {
	fmt.Println("plot '-' with points pointtype 6")

	for coord := range world {
		fmt.Printf("%d, %d\n", coord.x, coord.y)
	}
	
	fmt.Println("e")
}

func main() {
	// The world
	var world World
	world = make(World)
	
	// Define a starting world: the r-Pentomino
	world[Coord{1, 0}] = Cell{true, 0}
	world[Coord{0, 1}] = Cell{true, 0}
	world[Coord{1, 1}] = Cell{true, 0}
	world[Coord{1, 2}] = Cell{true, 0}
	world[Coord{2, 2}] = Cell{true, 0}
	
	gnuplotHeader()

	gnuplotWorld(world)
	
	for i := 0; i < 1000; i++ {
		world = world.Inflate().Tick().Deflate()
		gnuplotWorld(world)
	}
}
