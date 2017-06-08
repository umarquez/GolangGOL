package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"image/color"
	"strings"
)

const (
	// Default cell dimension
	CellSize = 8
	// Default screen width dimension
	ScreenWidth = 8 * 85
	// Default screen height dimension
	ScreenHeight = 8 * 85
)

var initialStatus Design
var cells []*Cell // Slice of cells to process

// FirstScene is de first and unique scene for this project
type FirstScene struct {
	Pattern Design
}

// SystemsList represents a collection of Systems
type SystemsList struct {
	RenderSystems *common.RenderSystem
	LifeSystems   *LifeSystem
}

// Type returns "FirstScene"
func (FirstScene) Type() string { return "FirstScene" }

// Preload pre-process assets an do actions before executes Setup
func (FirstScene) Preload() {}

// Setup generates the GOL cells grid, add the neighbors to each cell and draw
// the pattern to be processed.
func (scene FirstScene) Setup(w *ecs.World) {
	common.SetBackground(color.White)

	systems := SystemsList{
		&common.RenderSystem{},
		&LifeSystem{},
	}

	// Number of cells on screen
	maxCells := (ScreenWidth * ScreenHeight) / (CellSize * CellSize)
	// Number of cells per column
	maxCols := ScreenWidth / CellSize
	// Number of cells per row
	maxRows := ScreenHeight / CellSize

	// Generating cells
	for i := 0; i < maxCells; i++ {
		is := i * CellSize
		row := is / ScreenWidth
		col := (i + maxCols) % maxCols
		x := col * CellSize
		y := row * CellSize

		// Cell position on the screen
		cellPos := engo.Point{float32(x), float32(y)}

		// New cell on previous defined position
		cell := NewCell(cellPos)

		// Cell added to the cells list
		cells = append(cells, cell)

		// Check if current cell is an active cell on the
		position := engo.Point{float32(col), float32(row)}
		cell.status = IsActiveCell(scene.Pattern, position)

		systems.RenderSystems.AddByInterface(cell)
		systems.LifeSystems.Add(
			&cell.BasicEntity,
			&cell.RenderComponent,
			&cell.SpaceComponent,
			&cell.AliveComponent,
		)
	}

	// Adding neighbors to each cell if is possible.
	for i, cell := range cells {
		is := i * CellSize
		row := is / ScreenWidth
		col := (i + maxCols) % maxCols

		// fmt.Println(i, col, row)

		// W cell
		if col > 0 {
			cell.AddNeighbor(cells[i-1])
		}

		// N cell
		if row > 0 {
			cell.AddNeighbor(cells[i-maxCols])
		}

		// NW cell
		if col > 0 && row > 0 {
			cell.AddNeighbor(cells[i-maxCols-1])
		}

		// E cell
		if col < maxCols-1 {
			cell.AddNeighbor(cells[i+1])
		}

		// S cell
		if row < maxRows-1 {
			cell.AddNeighbor(cells[i+maxCols])
		}

		// SE cell
		if col < maxCols-1 && row < maxRows-1 {
			cell.AddNeighbor(cells[i+maxCols+1])
		}

		// NE cell
		if col < maxCols-1 && row > 0 {
			cell.AddNeighbor(cells[i-maxCols+1])
		}

		// SW cell
		if row < maxRows-1 && col > 0 {
			cell.AddNeighbor(cells[i+maxCols-1])
		}
	}

	w.AddSystem(systems.RenderSystems)
	w.AddSystem(systems.LifeSystems)

	fmt.Println("Cells ", len(cells))
}

func init() {
	/*
	 * ####################################################################
	 * # Six different ways to initialize the map.                        #
	 * # Each option use a custom tool build in the package that          #
	 * # satisfies the Design interface, that is used to set up each      #
	 * # cell according to the result of the querying for the status of   #
	 * # the (x, y) position at the scene's setup function.               #
	 * #                                                                  #
	 * # Comment the default and uncomment the option that you wanna test #
	 * # or create a new one.                                             #
	 * ####################################################################
	 */

	// Option #1: Simple glider in position 1, 1
	//
	// initialStatus = designs.NewSimpleGlider(engo.Point{1, 1})

	// Option #2: Gosper Glider Gun in de position 4, 4
	//
	// initialStatus = NewGosperGliderGun(engo.Point{4, 4})

	// Option #3: Simplex Noise generator with random values
	//
	// initialStatus = NewRandomSimplexNoise()

	// Option #4: A custom Simplex Noise generator that every time will
	// draw the same init status.
	//
	// initialStatus = &SimplexNoise{
	//	.7,
	//	.1,
	//	engo.Point{0,0},
	// }

	// Option #5: Alternate pattern [1 on/2 off]:
	// x..x..
	// .x..x.
	// ..x..x
	//
	// initialStatus = &AlternatePattern{}

	// Option #6: Tile pattern made of pulsars.
	// - First uses a free form to make one element.
	// - Then render the element on a canvas.
	// - Duplicates the canvas and flips each one: x axis, y axis and boths.
	// - Add the result patterns in a loop.

	pattern := NewFreeForm(
		engo.Point{2, 2},
	)
	pattern.Set(engo.Point{-2, -2}, true)
	pattern.Set(engo.Point{2, 0}, true)
	pattern.Set(engo.Point{3, 0}, true)
	pattern.Set(engo.Point{4, 0}, true)
	pattern.Set(engo.Point{0, 2}, true)
	pattern.Set(engo.Point{5, 2}, true)
	pattern.Set(engo.Point{0, 3}, true)
	pattern.Set(engo.Point{5, 3}, true)
	pattern.Set(engo.Point{0, 4}, true)
	pattern.Set(engo.Point{5, 4}, true)
	pattern.Set(engo.Point{2, 5}, true)
	pattern.Set(engo.Point{3, 5}, true)
	pattern.Set(engo.Point{4, 5}, true)

	design := NewCanvas(engo.Point{8, 8})
	design.RenderDesign(pattern)
	design.StartPoint = engo.Point{0, 0}

	a := *design
	b := *design
	c := *design
	d := *design

	b.FlipVertical()
	c.FlipHorizontal()
	d.FlipVertical()
	d.FlipHorizontal()

	mapContainer := NewMultiPattern()
	for x := 0; x < ScreenWidth; x += 17 {
		for y := 0; y < ScreenHeight; y += 17 {
			nw := a
			ne := b
			sw := c
			se := d
			nw.StartPoint = engo.Point{float32(x), float32(y)}
			ne.StartPoint = engo.Point{float32(x) + 9, float32(y)}
			sw.StartPoint = engo.Point{float32(x), float32(y) + 9}
			se.StartPoint = engo.Point{float32(x) + 9, float32(y) + 9}

			mapContainer.AddElement(&nw)
			mapContainer.AddElement(&ne)
			mapContainer.AddElement(&sw)
			mapContainer.AddElement(&se)
		}
	}

	initialStatus = mapContainer
}

func main() {
	consoleLeftMargin := strings.Repeat(" ", 20)
	fmt.Println(consoleLeftMargin + " _ _ _ _ _ ")
	fmt.Println(consoleLeftMargin + "|_|_|_|_|_|")
	fmt.Println(consoleLeftMargin + "|_|_|*|_|_|")
	fmt.Println(consoleLeftMargin + "|_|_|_|*|_|")
	fmt.Println(consoleLeftMargin + "|_|*|*|*|_|")
	fmt.Println(consoleLeftMargin + "|_|_|_|_|_|")
	fmt.Println("")
	fmt.Println("ng Conway's implementation using engo engine.")
	fmt.Println("Author: @umarquez")
	fmt.Println("")

	opts := engo.RunOptions{
		Title:      "ng ",
		Width:      ScreenWidth,  // 85 cols
		Height:     ScreenHeight, // 85 rows
		Fullscreen: false,
	}

	engo.Run(opts, FirstScene{initialStatus})
}
