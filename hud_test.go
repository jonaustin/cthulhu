package main

import (
	"testing"

	"game/engine"
	"game/render"
)

func TestStairsHint(t *testing.T) {
	if got := stairsHint(5, 5, 20, 20); got != "" {
		t.Fatalf("expected no hint far away, got %q", got)
	}
	if got := stairsHint(5, 5, 7, 6); got == "" {
		t.Fatal("expected hint within radius")
	}
	if got := stairsHint(5, 5, 5, 6); got == "" {
		t.Fatal("expected stronger hint adjacent")
	}
}

func TestBuildMiniMapMarksPlayerAndStairs(t *testing.T) {
	m := &engine.GameMap{
		Width:  5,
		Height: 5,
		Cells: [][]int{
			{engine.CellWall, engine.CellWall, engine.CellWall, engine.CellWall, engine.CellWall},
			{engine.CellWall, engine.CellEmpty, engine.CellEmpty, engine.CellEmpty, engine.CellWall},
			{engine.CellWall, engine.CellEmpty, engine.CellEmpty, engine.CellStairs, engine.CellWall},
			{engine.CellWall, engine.CellEmpty, engine.CellEmpty, engine.CellEmpty, engine.CellWall},
			{engine.CellWall, engine.CellWall, engine.CellWall, engine.CellWall, engine.CellWall},
		},
	}

	lines := buildMiniMap(m, 2, 2, 3, 2, 1)
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[1][1] != '@' {
		t.Fatalf("expected player '@' at center, got %q", lines[1][1])
	}
	if lines[1][2] != byte(render.StairsChar) {
		t.Fatalf("expected stairs %q at center-right, got %q", render.StairsChar, lines[1][2])
	}
}
