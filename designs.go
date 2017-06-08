package main

import (
	"engo.io/engo"
	"github.com/larspensjo/Go-simplex-noise/simplexnoise"
	"math/rand"
)

// Design is an interface used to work with any design
type Design interface {
	isActiveCell(position engo.Point) bool
}

// IsActiveCell returns the design's cell status
func IsActiveCell(d Design, position engo.Point) bool {
	return d.isActiveCell(position)
}

// GosperGliderGun Generates a Gosper Glider Gun
type GosperGliderGun struct {
	activeCells map[int]map[int]bool
	startPoint  engo.Point
}

// NewGosperGliderGun returns a *GosperGliderGun at startPoint
func NewGosperGliderGun(startPoint engo.Point) *GosperGliderGun {
	cells := make(map[int]map[int]bool)

	cells[0] = make(map[int]bool)
	cells[1] = make(map[int]bool)
	cells[10] = make(map[int]bool)
	cells[11] = make(map[int]bool)
	cells[12] = make(map[int]bool)
	cells[13] = make(map[int]bool)
	cells[14] = make(map[int]bool)
	cells[15] = make(map[int]bool)
	cells[16] = make(map[int]bool)
	cells[17] = make(map[int]bool)
	cells[20] = make(map[int]bool)
	cells[21] = make(map[int]bool)
	cells[22] = make(map[int]bool)
	cells[24] = make(map[int]bool)
	cells[34] = make(map[int]bool)
	cells[35] = make(map[int]bool)

	cells[0][4] = true
	cells[0][5] = true
	cells[1][4] = true
	cells[1][5] = true
	cells[10][4] = true
	cells[10][5] = true
	cells[10][6] = true
	cells[11][3] = true
	cells[11][7] = true
	cells[12][2] = true
	cells[12][8] = true
	cells[13][2] = true
	cells[13][8] = true
	cells[14][5] = true
	cells[15][3] = true
	cells[15][7] = true
	cells[16][4] = true
	cells[16][5] = true
	cells[16][6] = true
	cells[17][5] = true
	cells[20][2] = true
	cells[20][3] = true
	cells[20][4] = true
	cells[21][2] = true
	cells[21][3] = true
	cells[21][4] = true
	cells[22][1] = true
	cells[22][5] = true
	cells[24][0] = true
	cells[24][1] = true
	cells[24][5] = true
	cells[24][6] = true
	cells[34][2] = true
	cells[34][3] = true
	cells[35][2] = true
	cells[35][3] = true

	return &GosperGliderGun{
		activeCells: cells,
		startPoint:  startPoint,
	}
}

func (design *GosperGliderGun) isActiveCell(position engo.Point) bool {
	position.Subtract(design.startPoint)
	x := int(position.X)
	y := int(position.Y)

	return design.activeCells[x][y]
}

// SimplexNoise Generates a Simplex Noise Pattern
type SimplexNoise struct {
	Scale     float32
	Threshold float64
	Position  engo.Point
}

// SimplexNoise Generates a Simplex Noise Pattern
func NewRandomSimplexNoise() *SimplexNoise {
	return &SimplexNoise{
		Scale:     1,
		Threshold: 0.1,
		Position:  engo.Point{rand.Float32(), rand.Float32()},
	}
}

func (design *SimplexNoise) isActiveCell(position engo.Point) bool {
	position.Multiply(engo.Point{design.Scale, design.Scale})
	position.Add(design.Position)
	r := simplexnoise.Noise2(float64(position.X), float64(position.Y))

	return r > design.Threshold
}

// SimpleGlider Just builds a glider
type SimpleGlider struct {
	activeCells map[int]map[int]bool
	startPoint  engo.Point
}

// SimpleGlider Returns a SimpleGlider at startPoint
func NewSimpleGlider(startPoint engo.Point) *SimpleGlider {
	cells := make(map[int]map[int]bool)

	cells[0] = make(map[int]bool)
	cells[1] = make(map[int]bool)
	cells[2] = make(map[int]bool)
	cells[0][2] = true
	cells[1][0] = true
	cells[1][2] = true
	cells[2][1] = true
	cells[2][2] = true

	return &SimpleGlider{
		cells,
		startPoint,
	}
}

func (design *SimpleGlider) isActiveCell(position engo.Point) bool {
	position.Subtract(design.startPoint)
	x := int(position.X)
	y := int(position.Y)

	return design.activeCells[x][y]
}

// SimpleGlider Just builds a glider
type AlternatePattern struct{}

func (design *AlternatePattern) isActiveCell(position engo.Point) bool {
	r := (int(position.Y) + 3) % 3

	switch (int(position.X) + 3) % 3 {
	case 0:
		return r == 0
	case 1:
		return r == 1
	case 2:
		return r == 2
	}

	return false
}

// MultiPattern
type MultiPattern struct {
	elements []*Design
}

func NewMultiPattern() *MultiPattern {
	return &MultiPattern{
		[]*Design{},
	}
}

func (tpg *MultiPattern) AddElement(design Design) {
	tpg.elements = append(tpg.elements, &design)
}

func (tpg *MultiPattern) isActiveCell(position engo.Point) bool {
	result := false

	for _, element := range tpg.elements {
		e := *element
		result = result || IsActiveCell(e, position)
	}

	return result
}

type Square struct {
	activeCells map[int]map[int]bool
	startPoint  engo.Point
}

func NewSquare(startPoint engo.Point) *Square {
	cells := make(map[int]map[int]bool)

	cells[0] = make(map[int]bool)
	cells[1] = make(map[int]bool)
	cells[0][0] = true
	cells[0][1] = true
	cells[1][0] = true
	cells[1][1] = true

	return &Square{
		activeCells: cells,
		startPoint:  startPoint,
	}
}

func (design *Square) isActiveCell(position engo.Point) bool {
	position.Subtract(design.startPoint)
	x := int(position.X)
	y := int(position.Y)

	return design.activeCells[x][y]
}

type Canvas struct {
	size       engo.Point
	cells      map[int]map[int]bool
	StartPoint engo.Point
}

func NewCanvas(size engo.Point) *Canvas {
	cells := make(map[int]map[int]bool)

	for x := 0; x < int(size.X); x++ {
		cells[x] = make(map[int]bool)
		for y := 0; y < int(size.Y); y++ {
			cells[x][y] = false
		}
	}

	return &Canvas{
		cells: cells,
		size:  size,
	}
}

func (canvas *Canvas) RenderDesign(d Design) {
	for x := 0; x < int(canvas.size.X); x++ {
		for y := 0; y < int(canvas.size.Y); y++ {
			canvas.cells[x][y] = IsActiveCell(d, engo.Point{float32(x), float32(y)})
		}
	}
}

func (canvas *Canvas) FlipVertical() {
	cells := make(map[int]map[int]bool)

	for x := 0; x < int(canvas.size.X); x++ {
		newX := int(canvas.size.X) - 1 - x
		cells[newX] = make(map[int]bool)
		for y := 0; y < int(canvas.size.Y); y++ {
			cells[newX][y] = canvas.cells[x][y]
		}
	}

	canvas.cells = cells
}

func (canvas *Canvas) FlipHorizontal() {
	cells := make(map[int]map[int]bool)

	for y := 0; y < int(canvas.size.Y); y++ {
		newY := int(canvas.size.Y) - 1 - y
		for x := 0; x < int(canvas.size.X); x++ {
			if cells[x] == nil {
				cells[x] = make(map[int]bool)
			}
			cells[x][newY] = canvas.cells[x][y]
		}
	}

	canvas.cells = cells
}

func (canvas *Canvas) isActiveCell(position engo.Point) bool {
	position.Subtract(canvas.StartPoint)
	x := int(position.X)
	y := int(position.Y)

	return canvas.cells[x][y]
}

type FreeForm struct {
	cells      map[int]map[int]bool
	StartPoint engo.Point
}

func NewFreeForm(startPoint engo.Point) *FreeForm {
	return &FreeForm{
		cells:      make(map[int]map[int]bool),
		StartPoint: startPoint,
	}
}

func (freeForm *FreeForm) Set(position engo.Point, value bool) {
	if freeForm.cells[int(position.X)] == nil {
		freeForm.cells[int(position.X)] = make(map[int]bool)
	}

	freeForm.cells[int(position.X)][int(position.Y)] = value
}

func (freeForm *FreeForm) isActiveCell(position engo.Point) bool {
	position.Subtract(freeForm.StartPoint)
	x := int(position.X)
	y := int(position.Y)

	return freeForm.cells[x][y]
}
