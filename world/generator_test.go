package world

import (
	"testing"

	"game/engine"
)

func TestFloorGeneratorDeterministicWithSeed(t *testing.T) {
	g1 := NewFloorGenerator(32, 32, 1).WithSeed(123)
	m1 := g1.Generate()

	g2 := NewFloorGenerator(32, 32, 1).WithSeed(123)
	m2 := g2.Generate()

	if m1.Width != m2.Width || m1.Height != m2.Height {
		t.Fatalf("expected same dimensions, got %dx%d vs %dx%d", m1.Width, m1.Height, m2.Width, m2.Height)
	}

	for y := 0; y < m1.Height; y++ {
		for x := 0; x < m1.Width; x++ {
			if m1.Cells[y][x] != m2.Cells[y][x] {
				t.Fatalf("maps differ at (%d,%d): %d vs %d", x, y, m1.Cells[y][x], m2.Cells[y][x])
			}
		}
	}
}

func TestFloorGeneratorDifferentSeedsProduceDifferentMaps(t *testing.T) {
	g1 := NewFloorGenerator(32, 32, 1).WithSeed(1)
	m1 := g1.Generate()
	g2 := NewFloorGenerator(32, 32, 1).WithSeed(2)
	m2 := g2.Generate()

	diff := 0
	for y := 0; y < m1.Height; y++ {
		for x := 0; x < m1.Width; x++ {
			if m1.Cells[y][x] != m2.Cells[y][x] {
				diff++
			}
		}
	}
	if diff == 0 {
		t.Fatal("expected different seeds to produce different maps")
	}
}

func TestFloorGeneratorConnectivityAndStairsReachable(t *testing.T) {
	g := NewFloorGenerator(32, 32, 10).WithSeed(42)
	m := g.Generate()

	if !m.IsValid(g.SpawnPos.X, g.SpawnPos.Y) {
		t.Fatalf("spawn out of bounds: %+v", g.SpawnPos)
	}
	if !m.IsValid(g.StairsPos.X, g.StairsPos.Y) {
		t.Fatalf("stairs out of bounds: %+v", g.StairsPos)
	}

	if m.GetCell(g.SpawnPos.X, g.SpawnPos.Y) == engine.CellWall {
		t.Fatalf("spawn is on wall at %+v", g.SpawnPos)
	}
	if m.GetCell(g.StairsPos.X, g.StairsPos.Y) != engine.CellStairs {
		t.Fatalf("stairs cell not marked as stairs at %+v", g.StairsPos)
	}

	reachable, stairsReached := floodFillCount(m, g.SpawnPos)
	totalPassable := countPassable(m)

	if !stairsReached {
		t.Fatal("expected stairs to be reachable from spawn")
	}
	if reachable != totalPassable {
		t.Fatalf("expected all passable tiles to be connected: reachable=%d total=%d", reachable, totalPassable)
	}
}

func TestFloorGeneratorDepthAffectsOpenness(t *testing.T) {
	seed := int64(99)
	shallow := NewFloorGenerator(32, 32, 1).WithSeed(seed).Generate()
	deep := NewFloorGenerator(32, 32, 40).WithSeed(seed).Generate()

	shallowOpen := countPassable(shallow)
	deepOpen := countPassable(deep)

	if deepOpen >= shallowOpen {
		t.Fatalf("expected deeper floor to be less open: shallow=%d deep=%d", shallowOpen, deepOpen)
	}
}

func countPassable(m *engine.GameMap) int {
	total := 0
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if m.Cells[y][x] != engine.CellWall {
				total++
			}
		}
	}
	return total
}

func floodFillCount(m *engine.GameMap, start Point) (reachable int, stairsReached bool) {
	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}

	q := make([]Point, 0, m.Width*m.Height)
	head := 0

	q = append(q, start)
	visited[start.Y][start.X] = true

	for head < len(q) {
		cur := q[head]
		head++

		reachable++
		if m.Cells[cur.Y][cur.X] == engine.CellStairs {
			stairsReached = true
		}

		for _, d := range []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := cur.X+d.X, cur.Y+d.Y
			if !m.IsValid(nx, ny) {
				continue
			}
			if visited[ny][nx] {
				continue
			}
			if m.Cells[ny][nx] == engine.CellWall {
				continue
			}
			visited[ny][nx] = true
			q = append(q, Point{X: nx, Y: ny})
		}
	}

	return reachable, stairsReached
}
