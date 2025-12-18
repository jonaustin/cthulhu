package engine

import (
	"math"
	"testing"
)

func TestRaycasterCastRay(t *testing.T) {
	r := NewRaycaster(120, 40)
	m := NewTestMap()
	// Player at center of map facing right (east)
	p := NewPlayer(8.5, 8.5, 0)

	// Cast ray straight ahead - should hit east wall eventually
	dist := r.castRay(p, m, 0)
	if dist <= 0 {
		t.Errorf("Expected positive distance, got %f", dist)
	}
	if dist > r.MaxDist {
		t.Errorf("Distance %f exceeds max distance %f", dist, r.MaxDist)
	}
}

func TestRaycasterCastRayWithStairsTracksNearestStairsBeforeWall(t *testing.T) {
	r := NewRaycaster(120, 40)
	m := &GameMap{
		Width:  5,
		Height: 5,
		Cells: [][]int{
			{CellWall, CellWall, CellWall, CellWall, CellWall},
			{CellWall, CellEmpty, CellEmpty, CellEmpty, CellWall},
			{CellWall, CellEmpty, CellEmpty, CellStairs, CellWall},
			{CellWall, CellEmpty, CellEmpty, CellEmpty, CellWall},
			{CellWall, CellWall, CellWall, CellWall, CellWall},
		},
	}

	p := NewPlayer(1.5, 2.5, 0) // facing east
	wallDist, stairsDist := r.castRayWithStairs(p, m, 0)

	if wallDist <= 0 {
		t.Fatalf("expected wallDist > 0, got %f", wallDist)
	}
	if !(stairsDist > 0) {
		t.Fatalf("expected stairsDist > 0, got %f", stairsDist)
	}
	if stairsDist >= wallDist {
		t.Fatalf("expected stairsDist < wallDist, got stairs=%f wall=%f", stairsDist, wallDist)
	}
}

func TestRaycasterFishEyeCorrection(t *testing.T) {
	r := NewRaycaster(120, 40)
	m := NewTestMap()
	// Player facing a wall head-on
	p := NewPlayer(8.5, 2.5, math.Pi/2) // facing south

	// Cast rays at center and edges
	centerDist := r.castRay(p, m, p.Angle)
	leftDist := r.castRay(p, m, p.Angle-r.FOV/2)
	rightDist := r.castRay(p, m, p.Angle+r.FOV/2)

	// All rays should hit walls at some positive distance
	if centerDist <= 0 || leftDist <= 0 || rightDist <= 0 {
		t.Errorf("Rays should hit walls: center=%f left=%f right=%f",
			centerDist, leftDist, rightDist)
	}
}

func TestRaycasterSetScreenSize(t *testing.T) {
	r := NewRaycaster(120, 40)
	if r.ScreenWidth != 120 || r.ScreenHeight != 40 {
		t.Errorf("Initial size wrong: got %dx%d", r.ScreenWidth, r.ScreenHeight)
	}

	r.SetScreenSize(80, 24)
	if r.ScreenWidth != 80 || r.ScreenHeight != 24 {
		t.Errorf("After resize wrong: got %dx%d", r.ScreenWidth, r.ScreenHeight)
	}
}

func TestNewPlayer(t *testing.T) {
	p := NewPlayer(5.5, 3.5, math.Pi/4)
	if p.X != 5.5 || p.Y != 3.5 {
		t.Errorf("Position wrong: got (%f, %f)", p.X, p.Y)
	}
	if p.Angle != math.Pi/4 {
		t.Errorf("Angle wrong: got %f", p.Angle)
	}
}

func TestPlayerDirection(t *testing.T) {
	// Player facing east (angle 0)
	p := NewPlayer(0, 0, 0)
	if math.Abs(p.DirX()-1.0) > 0.001 {
		t.Errorf("DirX for angle 0 should be 1, got %f", p.DirX())
	}
	if math.Abs(p.DirY()) > 0.001 {
		t.Errorf("DirY for angle 0 should be 0, got %f", p.DirY())
	}

	// Player facing south (angle π/2)
	p2 := NewPlayer(0, 0, math.Pi/2)
	if math.Abs(p2.DirX()) > 0.001 {
		t.Errorf("DirX for angle π/2 should be 0, got %f", p2.DirX())
	}
	if math.Abs(p2.DirY()-1.0) > 0.001 {
		t.Errorf("DirY for angle π/2 should be 1, got %f", p2.DirY())
	}
}
