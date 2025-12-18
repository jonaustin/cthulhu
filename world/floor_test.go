package world

import (
	"testing"

	"game/engine"
)

func TestFloorManagerGenerateFirstFloor(t *testing.T) {
	fm := NewFloorManager()
	fm.Generator.WithSeed(123)

	floor := fm.GenerateFirstFloor()
	if floor.Depth != 1 {
		t.Fatalf("expected depth 1, got %d", floor.Depth)
	}
	if floor.Map == nil {
		t.Fatal("expected generated map")
	}
	if floor.Map.GetCell(floor.SpawnPos.X, floor.SpawnPos.Y) == engine.CellWall {
		t.Fatalf("spawn is on wall at %+v", floor.SpawnPos)
	}
	if floor.Map.GetCell(floor.StairsPos.X, floor.StairsPos.Y) != engine.CellStairs {
		t.Fatalf("stairs not present at %+v", floor.StairsPos)
	}
	if fm.GetCurrentDepth() != 1 {
		t.Fatalf("expected GetCurrentDepth()=1, got %d", fm.GetCurrentDepth())
	}
}

func TestFloorManagerDescendToNextFloorIncrementsDepth(t *testing.T) {
	fm := NewFloorManager()
	fm.Generator.WithSeed(321)

	f1 := fm.GenerateFirstFloor()
	f2 := fm.DescendToNextFloor()

	if f2.Depth != f1.Depth+1 {
		t.Fatalf("expected depth %d, got %d", f1.Depth+1, f2.Depth)
	}
	if f1.Map == f2.Map {
		t.Fatal("expected a new map instance on descent")
	}
	if fm.GetCurrentDepth() != f2.Depth {
		t.Fatalf("expected GetCurrentDepth()=%d, got %d", f2.Depth, fm.GetCurrentDepth())
	}
}
