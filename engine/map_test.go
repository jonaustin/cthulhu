package engine

import "testing"

func TestNewTestMap(t *testing.T) {
	m := NewTestMap()

	if m.Width != 16 || m.Height != 16 {
		t.Errorf("expected 16x16 map, got %dx%d", m.Width, m.Height)
	}

	// Perimeter should be all walls
	for x := 0; x < 16; x++ {
		if !m.IsWall(x, 0) {
			t.Errorf("top perimeter should be wall at x=%d", x)
		}
		if !m.IsWall(x, 15) {
			t.Errorf("bottom perimeter should be wall at x=%d", x)
		}
	}
	for y := 0; y < 16; y++ {
		if !m.IsWall(0, y) {
			t.Errorf("left perimeter should be wall at y=%d", y)
		}
		if !m.IsWall(15, y) {
			t.Errorf("right perimeter should be wall at y=%d", y)
		}
	}

	// Open area should be walkable (spawn point at 8,6 - verified open)
	if m.IsWall(8, 6) {
		t.Error("spawn area (8,6) should be walkable")
	}

	// Check an internal wall exists (from the L-shaped wall)
	if !m.IsWall(3, 3) {
		t.Error("internal wall at (3,3) should exist")
	}
}

func TestIsValid(t *testing.T) {
	m := NewTestMap()

	// Valid coordinates
	if !m.IsValid(0, 0) {
		t.Error("(0,0) should be valid")
	}
	if !m.IsValid(15, 15) {
		t.Error("(15,15) should be valid")
	}
	if !m.IsValid(8, 8) {
		t.Error("(8,8) should be valid")
	}

	// Invalid coordinates
	if m.IsValid(-1, 0) {
		t.Error("(-1,0) should be invalid")
	}
	if m.IsValid(16, 0) {
		t.Error("(16,0) should be invalid")
	}
	if m.IsValid(0, -1) {
		t.Error("(0,-1) should be invalid")
	}
	if m.IsValid(0, 16) {
		t.Error("(0,16) should be invalid")
	}
}

func TestGetCell(t *testing.T) {
	m := NewTestMap()

	// Out of bounds returns wall
	if m.GetCell(-1, 0) != CellWall {
		t.Error("out of bounds should return CellWall")
	}

	// Valid wall
	if m.GetCell(0, 0) != CellWall {
		t.Error("corner should be CellWall")
	}

	// Valid empty
	if m.GetCell(8, 6) != CellEmpty {
		t.Error("spawn area (8,6) should be CellEmpty")
	}
}
