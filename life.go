package main

import (
	"engo.io/ecs"
	"engo.io/engo/common"
)

const (
	LifeStepInterval float32 = 0.1
)

// AliveComponent represents each cell's dot
type AliveComponent struct {
	neighbors  []*Cell
	status     bool
	lastStatus bool
	updated    bool
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// updateStatus just update the status of each cell based on the number of
// active neighbors.
func (c *AliveComponent) updateStatus() {
	total := 0
	for _, neighbor := range c.neighbors {
		if neighbor.updated && neighbor.lastStatus {
			total += 1
		} else if !neighbor.updated && neighbor.status {
			total += 1
		}
	}

	/*
	 * ##################################
	 * # Game Of Life Rules evaluation. #
	 * ##################################
	 */
	c.lastStatus = c.status

	if c.status {
		// Any live cell with fewer than two live neighbours dies, as if caused
		// by underpopulation.
		if total < 2 {
			c.status = false
		}

		// Any live cell with two or three live neighbours lives on to the next
		// generation.
		// if total == 2 || total == 3 { s.status = true }

		// Any live cell with more than three live neighbours dies, as if by
		// overpopulation.
		if total > 3 {
			c.status = false
		}
	} else {
		// Any dead cell with exactly three live neighbours becomes a live
		// cell, as if by reproduction.
		if total == 3 {
			c.status = true
		}
	}

	c.updated = true
}

// AliveEntity is a simple cell
type AliveEntity struct {
	*ecs.BasicEntity
	*common.RenderComponent
	*common.SpaceComponent
	*AliveComponent
}

// LifeSystem is used to manage each cell's status on each step
type LifeSystem struct {
	since       float32
	gen         int
	isFirstTime bool
	world       *ecs.World
	entities    []AliveEntity
}

// New initialize the system
func (l *LifeSystem) New(w *ecs.World) {
	l.world = w
	l.isFirstTime = true
}

// Add append an entity to the system
func (l *LifeSystem) Add(
	basic *ecs.BasicEntity,
	render *common.RenderComponent,
	space *common.SpaceComponent,
	alive *AliveComponent,
) {
	l.entities = append(l.entities, AliveEntity{
		basic,
		render,
		space,
		alive,
	})
}

// Remove deletes an entity
func (l *LifeSystem) Remove(basic ecs.BasicEntity) {
	deleteIndex := -1
	for index, e := range l.entities {
		if e.BasicEntity.ID() == basic.ID() {
			deleteIndex = index
			break
		}
	}
	if deleteIndex >= 0 {
		l.entities = append(l.entities[:deleteIndex], l.entities[deleteIndex+1:]...)
	}
}

// clearEntitiesStatus set each cell's updated status to false
func (l *LifeSystem) clearEntitiesStatus() {
	for _, entity := range l.entities {
		entity.updated = false
	}
}

// Update process each cell to get the new status and shows or hide the dot
func (l *LifeSystem) Update(dt float32) {
	l.since += dt

	if l.since > LifeStepInterval || l.isFirstTime {
		l.gen += 1
		l.clearEntitiesStatus()
		var sys *common.RenderSystem

		for _, system := range l.world.Systems() {
			switch system.(type) {
			case *common.RenderSystem:
				sys = system.(*common.RenderSystem)
				break
			}
		}

		for _, entity := range l.entities {
			if !l.isFirstTime {
				entity.updateStatus()
			}

			if entity.status {
				sys.AddByInterface(entity.AliveComponent)
			} else {
				sys.Remove(entity.AliveComponent.BasicEntity)
			}
		}
		l.since = 0
	}

	if l.isFirstTime {
		l.isFirstTime = false
	}
}
