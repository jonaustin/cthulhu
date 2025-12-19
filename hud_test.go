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

func TestBuildMiniMapRectRespectsDifferentRadii(t *testing.T) {
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

	lines := buildMiniMapRect(m, 2, 2, 3, 2, 2, 1) // 5 wide, 3 tall
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if got := len([]rune(lines[0])); got != 5 {
		t.Fatalf("expected width 5, got %d", got)
	}
}

func TestMiniMapStartXLeavesRightMargin(t *testing.T) {
	if got := miniMapStartX(80, 10); got != 69 {
		t.Fatalf("expected startX 69, got %d", got)
	}
	if got := miniMapStartX(5, 10); got != 0 {
		t.Fatalf("expected startX clamped to 0, got %d", got)
	}
}
