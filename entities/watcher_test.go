package entities

import (
	"math"
	"testing"
)

func TestWatcherManagerCountsByDepth(t *testing.T) {
	fov := math.Pi / 3
	seed := int64(123)

	cases := []struct {
		depth      int
		min, max   int
		expectNone bool
	}{
		{depth: 1, min: 0, max: 0, expectNone: true},
		{depth: WatcherStartDepth, min: 1, max: 2},
		{depth: watcherDepthTier1Max, min: 1, max: 2},
		{depth: watcherDepthTier1Max + 1, min: 2, max: 3},
		{depth: watcherDepthTier2Max, min: 2, max: 3},
		{depth: watcherDepthTier2Max + 1, min: 3, max: 4},
	}

	for _, tc := range cases {
		wm := NewWatcherManager(tc.depth, seed, fov)
		got := len(wm.Watchers)
		if tc.expectNone {
			if got != 0 {
				t.Fatalf("depth %d: expected 0 watchers, got %d", tc.depth, got)
			}
			continue
		}
		if got < tc.min || got > tc.max {
			t.Fatalf("depth %d: expected %d-%d watchers, got %d", tc.depth, tc.min, tc.max, got)
		}
	}
}

func TestWatcherManagerRanges(t *testing.T) {
	fov := math.Pi / 3
	wm := NewWatcherManager(20, 42, fov)

	minEdge, maxEdge := edgeOffsetRange(fov)
	expectedMin := (fov * 0.5) * (1.0 - WatcherEdgeThreshold)
	if math.Abs(minEdge-expectedMin) > 1e-9 {
		t.Fatalf("unexpected min edge: got %f, want %f", minEdge, expectedMin)
	}
	if math.Abs(maxEdge-(fov*0.5)) > 1e-9 {
		t.Fatalf("unexpected max edge: got %f, want %f", maxEdge, fov*0.5)
	}
	for _, w := range wm.Watchers {
		if w.Angle < minEdge || w.Angle > maxEdge {
			t.Fatalf("angle out of range: got %f, want %f-%f", w.Angle, minEdge, maxEdge)
		}
		if w.Distance < WatcherMinDistance || w.Distance > WatcherMaxDistance {
			t.Fatalf("distance out of range: got %f, want %f-%f", w.Distance, WatcherMinDistance, WatcherMaxDistance)
		}
		if w.Side != -1 && w.Side != 1 {
			t.Fatalf("unexpected side %d", w.Side)
		}
		if w.Drift != -1 && w.Drift != 1 {
			t.Fatalf("unexpected drift %f", w.Drift)
		}
	}
}

func TestWatcherManagerUpdateKeepsEdgeBand(t *testing.T) {
	fov := math.Pi / 3
	wm := NewWatcherManager(30, 7, fov)
	minEdge, maxEdge := edgeOffsetRange(fov)

	for i := 0; i < 500; i++ {
		wm.Update()
		for _, w := range wm.Watchers {
			if w.Angle < minEdge || w.Angle > maxEdge {
				t.Fatalf("tick %d: angle out of range %f (min %f max %f)", i, w.Angle, minEdge, maxEdge)
			}
			if w.Side != -1 && w.Side != 1 {
				t.Fatalf("tick %d: unexpected side %d", i, w.Side)
			}
		}
	}
}

func TestWatcherSpritesMatchVisibleCount(t *testing.T) {
	fov := math.Pi / 3
	wm := NewWatcherManager(20, 99, fov)
	wm.Update()

	visible := wm.VisibleCount()
	sprites := wm.Sprites(120, 40)
	if len(sprites) != visible {
		t.Fatalf("expected %d sprites, got %d", visible, len(sprites))
	}
}
