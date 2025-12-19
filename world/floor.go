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
	MapWidth     int
	MapHeight    int
}

const (
	DefaultMapWidth  = 32
	DefaultMapHeight = 32
)

func NewFloorManager() *FloorManager {
	return NewFloorManagerWithSize(DefaultMapWidth, DefaultMapHeight)
}

func NewFloorManagerWithSize(width, height int) *FloorManager {
	if width <= 0 {
		width = DefaultMapWidth
	}
	if height <= 0 {
		height = DefaultMapHeight
	}

	baseSeed := time.Now().UnixNano()
	return &FloorManager{
		MapWidth:  width,
		MapHeight: height,
		Generator: NewFloorGenerator(width, height, 1).WithSeed(baseSeed),
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

func (fm *FloorManager) TeleportToDepth(depth int) *Floor {
	if depth < 1 {
		depth = 1
	}
	return fm.generateAtDepth(depth)
}

func (fm *FloorManager) GetCurrentDepth() int {
	if fm.CurrentFloor == nil {
		return 0
	}
	return fm.CurrentFloor.Depth
}

func (fm *FloorManager) generateAtDepth(depth int) *Floor {
	if fm.Generator == nil {
		w := fm.MapWidth
		h := fm.MapHeight
		if w <= 0 {
			w = DefaultMapWidth
		}
		if h <= 0 {
			h = DefaultMapHeight
		}
		fm.Generator = NewFloorGenerator(w, h, depth)
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
