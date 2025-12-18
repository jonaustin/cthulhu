package world

import "testing"

func TestCorruptionCalculateLevel(t *testing.T) {
	c := NewCorruption()

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
		if got != tc.want {
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
