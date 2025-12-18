package engine

import (
	"math"
	"testing"
)

func TestNewPlayerAtCell(t *testing.T) {
	p := NewPlayerAtCell(3, 4, math.Pi)
	if p.X != 3.5 || p.Y != 4.5 {
		t.Errorf("expected centered position (3.5,4.5), got (%f,%f)", p.X, p.Y)
	}
	if p.Angle != math.Pi {
		t.Errorf("expected angle π, got %f", p.Angle)
	}
}

func TestPlayerRotateLeftRight(t *testing.T) {
	p := NewPlayer(0, 0, 0)

	p.RotateRight()
	if math.Abs(p.Angle-math.Pi/2) > 0.00001 {
		t.Errorf("RotateRight from 0 should be π/2, got %f", p.Angle)
	}

	p.RotateRight()
	if math.Abs(p.Angle-math.Pi) > 0.00001 {
		t.Errorf("RotateRight from π/2 should be π, got %f", p.Angle)
	}

	p.RotateLeft()
	if math.Abs(p.Angle-math.Pi/2) > 0.00001 {
		t.Errorf("RotateLeft from π should be π/2, got %f", p.Angle)
	}

	p2 := NewPlayer(0, 0, 0)
	p2.RotateLeft()
	if math.Abs(p2.Angle-3*math.Pi/2) > 0.00001 {
		t.Errorf("RotateLeft from 0 should wrap to 3π/2, got %f", p2.Angle)
	}
}

func TestPlayerMoveForwardBackward(t *testing.T) {
	m := NewTestMap()

	p := NewPlayerAtCell(8, 6, 0) // east
	p.MoveForward(m)
	if p.X != 9.5 || p.Y != 6.5 {
		t.Errorf("expected move east to (9.5,6.5), got (%f,%f)", p.X, p.Y)
	}

	p.MoveBackward(m)
	if p.X != 8.5 || p.Y != 6.5 {
		t.Errorf("expected move back west to (8.5,6.5), got (%f,%f)", p.X, p.Y)
	}
}

func TestPlayerMovementBlockedByWall(t *testing.T) {
	m := NewTestMap()

	// Cell (1,1) has a wall at (0,1) to the west (perimeter).
	p := NewPlayerAtCell(1, 1, math.Pi) // west
	p.MoveForward(m)
	if p.X != 1.5 || p.Y != 1.5 {
		t.Errorf("expected blocked by wall, got (%f,%f)", p.X, p.Y)
	}

	// Backward should move east into (2,1).
	p.MoveBackward(m)
	if p.X != 2.5 || p.Y != 1.5 {
		t.Errorf("expected move backward to (2.5,1.5), got (%f,%f)", p.X, p.Y)
	}
}

func TestPlayerMoveForwardSnapsAngleAndNormalizes(t *testing.T) {
	m := NewTestMap()

	// -π/2 is north (up) in this coordinate system.
	p := NewPlayerAtCell(8, 6, -math.Pi/2)
	p.MoveForward(m)
	if p.X != 8.5 || p.Y != 5.5 {
		t.Errorf("expected move north to (8.5,5.5), got (%f,%f)", p.X, p.Y)
	}
}
