package main

import (
	"testing"

	"game/engine"
	"game/render"
	"game/world"

	"github.com/gdamore/tcell/v2"
)

func newTestGameForCheats(t *testing.T) *Game {
	t.Helper()

	fm := world.NewFloorManager()
	fm.Generator.WithSeed(123)
	floor := fm.GenerateFirstFloor()

	return &Game{
		Running:      true,
		ShowMiniMap:  true,
		ShowWatchers: true,
		CorruptState: world.NewCorruption(),
		FloorManager: fm,
		Floor:        floor,
		GameMap:      floor.Map,
		Player:       engine.NewPlayerAtCell(floor.SpawnPos.X, floor.SpawnPos.Y, 0),
	}
}

func TestCheatMenuToggleMap(t *testing.T) {
	g := newTestGameForCheats(t)
	g.openCheatMenu()

	if !g.ShowMiniMap {
		t.Fatal("expected minimap on by default")
	}

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, 'm', tcell.ModNone))
	if g.ShowMiniMap {
		t.Fatal("expected minimap toggled off")
	}
}

func TestCheatMenuToggleWatchers(t *testing.T) {
	g := newTestGameForCheats(t)
	g.openCheatMenu()

	if !g.ShowWatchers {
		t.Fatal("expected watchers on by default")
	}

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, 'w', tcell.ModNone))
	if g.ShowWatchers {
		t.Fatal("expected watchers toggled off")
	}

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, 'w', tcell.ModNone))
	if !g.ShowWatchers {
		t.Fatal("expected watchers toggled on")
	}
}

func TestCheatMenuTeleportToDepth(t *testing.T) {
	g := newTestGameForCheats(t)
	g.openCheatMenu()

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, 't', tcell.ModNone))
	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, '5', tcell.ModNone))
	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone))

	if g.Floor == nil || g.Floor.Depth != 5 {
		t.Fatalf("expected depth 5 after teleport, got %+v", g.Floor)
	}
	if g.cheatMode != cheatModeMain {
		t.Fatalf("expected cheat mode main after teleport, got %v", g.cheatMode)
	}
}

func TestCheatMenuTuneAdjustsVisualConfig(t *testing.T) {
	render.ResetVisualConfig()
	defer render.ResetVisualConfig()

	g := newTestGameForCheats(t)
	g.openCheatMenu()

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, 'v', tcell.ModNone))
	before := render.GetVisualConfig()

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone))
	after := render.GetVisualConfig()
	if after.VisualScale <= before.VisualScale {
		t.Fatalf("expected visual scale to increase, got %f <= %f", after.VisualScale, before.VisualScale)
	}

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone))
	before = render.GetVisualConfig()
	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone))
	after = render.GetVisualConfig()
	if after.MaxCharGlitchChance >= before.MaxCharGlitchChance {
		t.Fatalf("expected char glitch chance to decrease, got %f >= %f", after.MaxCharGlitchChance, before.MaxCharGlitchChance)
	}
}

func TestCheatMenuAdjustCorruptionBias(t *testing.T) {
	g := newTestGameForCheats(t)
	g.openCheatMenu()

	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, '+', tcell.ModNone))
	if got := g.CorruptState.GetBias(); got <= 0 {
		t.Fatalf("expected positive bias after '+', got %f", got)
	}

	before := g.CorruptState.GetBias()
	g.handleCheatEvent(tcell.NewEventKey(tcell.KeyRune, '-', tcell.ModNone))
	if got := g.CorruptState.GetBias(); got >= before {
		t.Fatalf("expected bias to decrease after '-', got %f (before %f)", got, before)
	}
}

func TestCheatMenuEscapeDoesNotQuitGame(t *testing.T) {
	g := newTestGameForCheats(t)
	g.openCheatMenu()

	g.processEvent(tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone))

	if !g.Running {
		t.Fatal("expected game to keep running when closing cheat menu")
	}
	if g.cheatMenuOpen {
		t.Fatal("expected cheat menu to be closed")
	}
}
