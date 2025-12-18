package world

import (
	"time"

	"game/engine"
)

type Floor struct {
	Map       *engine.GameMap
	Depth     int
	SpawnPos  Point
	StairsPos Point
}

type FloorManager struct {
	CurrentFloor *Floor
	Generator    *FloorGenerator
}

const (
	defaultMapWidth  = 32
	defaultMapHeight = 32
)

func NewFloorManager() *FloorManager {
	baseSeed := time.Now().UnixNano()
	return &FloorManager{
		Generator: NewFloorGenerator(defaultMapWidth, defaultMapHeight, 1).WithSeed(baseSeed),
	}
}

func (fm *FloorManager) GenerateFirstFloor() *Floor {
	return fm.generateAtDepth(1)
}

func (fm *FloorManager) DescendToNextFloor() *Floor {
	nextDepth := 1
	if fm.CurrentFloor != nil {
		nextDepth = fm.CurrentFloor.Depth + 1
	}
	return fm.generateAtDepth(nextDepth)
}

func (fm *FloorManager) GetCurrentDepth() int {
	if fm.CurrentFloor == nil {
		return 0
	}
	return fm.CurrentFloor.Depth
}

func (fm *FloorManager) generateAtDepth(depth int) *Floor {
	if fm.Generator == nil {
		fm.Generator = NewFloorGenerator(defaultMapWidth, defaultMapHeight, depth)
	}

	fm.Generator.Depth = depth
	m := fm.Generator.Generate()

	f := &Floor{
		Map:       m,
		Depth:     depth,
		SpawnPos:  fm.Generator.SpawnPos,
		StairsPos: fm.Generator.StairsPos,
	}
	fm.CurrentFloor = f
	return f
}
