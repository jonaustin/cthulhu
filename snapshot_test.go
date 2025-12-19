package main

import (
	"math"
	"os"
	"path/filepath"
	"testing"
	"time"

	"game/engine"
	"game/render"
	"game/world"

	"github.com/gdamore/tcell/v2"
)

func TestSnapshotRenderToFile(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init screen: %v", err)
	}
	defer screen.Fini()

	const (
		width  = 64
		height = 24
	)
	screen.SetSize(width, height)

	fm := world.NewFloorManagerWithSize(20, 20)
	fm.Generator.WithSeed(123)
	floor := fm.TeleportToDepth(15)
	if floor == nil || floor.Map == nil {
		t.Fatal("expected generated floor")
	}

	player := engine.NewPlayerAtCell(floor.SpawnPos.X, floor.SpawnPos.Y, -math.Pi/2)
	raycaster := engine.NewRaycaster(width, height)

	corruption := world.NewCorruption()
	corruption.Update(floor.Depth)
	effects := render.NewEffectsContext(corruption.Depth, corruption.GetLevel(), corruption.Ticks)

	raycaster.RenderWithEffects(screen, player, floor.Map, effects, floor.Watchers)
	screen.Show()

	lines := captureScreenLines(screen, width, height)
	if len(lines) != height {
		t.Fatalf("expected %d lines, got %d", height, len(lines))
	}

	meta := snapshotMeta{
		Depth:      floor.Depth,
		Corruption: corruption.GetLevel(),
		Ticks:      corruption.Ticks,
		Width:      width,
		Height:     height,
		Timestamp:  time.Unix(0, 0).UTC(),
	}
	path := filepath.Join(t.TempDir(), snapshotFilename(meta))
	if err := writeSnapshotFile(path, meta, lines); err != nil {
		t.Fatalf("write snapshot: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read snapshot: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected snapshot data")
	}
}
