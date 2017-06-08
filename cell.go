package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"image/color"
	"math/rand"
)

// Cell is a square on the screen.
type Cell struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	AliveComponent
}

// NewCell creates a cell on a given position.
func NewCell(position engo.Point) *Cell {
	aliveSize := float32(CellSize * 0.5)
	aliveMargin := (CellSize - aliveSize) / 2

	alivePos := engo.Point{aliveMargin, aliveMargin}
	alivePos.Add(position)

	return &Cell{
		ecs.NewBasic(),
		common.RenderComponent{
			Drawable: common.Rectangle{
				BorderWidth: 2,
				BorderColor: color.RGBA{180, 180, 180, 255},
			},
			Color: color.Transparent,
		},
		common.SpaceComponent{
			Position: position,
			Width:    CellSize,
			Height:   CellSize,
		},
		AliveComponent{
			[]*Cell{},
			false,
			false,
			false,
			ecs.NewBasic(),
			common.RenderComponent{
				Drawable: common.Circle{},
				Color:    color.RGBA{20, 20, 20, 255},
			},
			common.SpaceComponent{
				Position: alivePos,
				Width:    aliveSize,
				Height:   aliveSize,
			},
		},
	}
}

// AddNeighbor add the given cell as a neighbor.
func (c *Cell) AddNeighbor(n *Cell) {
	c.neighbors = append(c.neighbors, n)
}

// RandomActivate set the cell status according to the given provability.
func (c *Cell) RandomActivate(provability float32) bool {
	r := c.status

	if rand.Float32() < provability {
		r = true
	}
	c.status = r

	return r
}
