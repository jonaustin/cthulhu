package world

import (
	"math"
	"testing"
)

func TestCorruptionCalculateLevel(t *testing.T) {
	c := NewCorruption()

	const eps = 1e-9
	cases := []struct {
		depth int
		want  float64
	}{
		{depth: 1, want: 0.0},
		{depth: 9, want: 0.0},
		{depth: 10, want: 0.0},
		{depth: 11, want: 0.025},
		{depth: 50, want: 1.0},
		{depth: 60, want: 1.0},
	}

	for _, tc := range cases {
		got := c.calculateLevel(tc.depth)
		if math.Abs(got-tc.want) > eps {
			t.Fatalf("depth %d: expected %f, got %f", tc.depth, tc.want, got)
		}
	}
}

func TestCorruptionUpdateSetsStateAndTicks(t *testing.T) {
	c := NewCorruption()
	if c.Ticks != 0 {
		t.Fatalf("expected initial ticks 0, got %d", c.Ticks)
	}

	c.Update(1)
	if c.Depth != 1 {
		t.Fatalf("expected depth 1, got %d", c.Depth)
	}
	if c.Level != 0.0 {
		t.Fatalf("expected level 0.0 at depth 1, got %f", c.Level)
	}
	if c.Ticks != 1 {
		t.Fatalf("expected ticks 1, got %d", c.Ticks)
	}

	c.Update(20)
	if c.Depth != 20 {
		t.Fatalf("expected depth 20, got %d", c.Depth)
	}
	if c.GetLevel() <= 0.0 {
		t.Fatalf("expected level > 0.0 at depth 20, got %f", c.GetLevel())
	}
	if c.Ticks != 2 {
		t.Fatalf("expected ticks 2, got %d", c.Ticks)
	}
}

func TestCorruptionBiasAdjustsGetLevel(t *testing.T) {
	c := NewCorruption()

	c.Update(1) // base = 0.0
	c.AdjustBias(0.5)
	if got := c.GetLevel(); got != 0.5 {
		t.Fatalf("expected level 0.5 with bias at depth 1, got %f", got)
	}

	c.AdjustBias(-0.75) // bias now -0.25
	if got := c.GetLevel(); got != 0.0 {
		t.Fatalf("expected clamped level 0.0 with negative bias, got %f", got)
	}

	c.AdjustBias(2.0) // clamp bias to 1
	c.Update(60)      // base = 1.0
	if got := c.GetLevel(); got != 1.0 {
		t.Fatalf("expected clamped level 1.0 at max, got %f", got)
	}
	if got := c.GetBias(); got != 1.0 {
		t.Fatalf("expected bias clamped to 1.0, got %f", got)
	}
}
