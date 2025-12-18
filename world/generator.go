package world

import (
	"math"
	"math/rand"
	"time"

	"game/engine"
)

type Point struct {
	X, Y int
}

// FloorGenerator generates a new map for a given depth.
//
// It uses a connected "drunk walk" carve so all empty tiles are reachable
// from SpawnPos by construction, then places StairsPos at the farthest
// reachable tile.
type FloorGenerator struct {
	Width, Height int
	Depth         int
	Seed          int64

	SpawnPos  Point
	StairsPos Point
}

const (
	minMapSize          = 5
	depthScaleMax       = 50.0
	seedDepthMultiplier = int64(1_000_003)
	baseOpenFraction    = 0.50
	openFractionDrop    = 0.22
	minOpenFraction     = 0.28
	baseTurnChance      = 0.25
	turnChanceIncrease  = 0.45
	maxStepsPerCell     = 12
	minTargetOpenCells  = 2
)

func NewFloorGenerator(width, height, depth int) *FloorGenerator {
	return &FloorGenerator{
		Width:  width,
		Height: height,
		Depth:  depth,
		Seed:   time.Now().UnixNano(),
	}
}

func (g *FloorGenerator) WithSeed(seed int64) *FloorGenerator {
	g.Seed = seed
	return g
}

// Generate returns a newly generated map and updates SpawnPos/StairsPos.
func (g *FloorGenerator) Generate() *engine.GameMap {
	w, h := g.Width, g.Height
	if w < minMapSize {
		w = minMapSize
	}
	if h < minMapSize {
		h = minMapSize
	}

	rng := rand.New(rand.NewSource(g.Seed + int64(g.Depth)*seedDepthMultiplier))

	m := newSolidWallMap(w, h)

	spawn := Point{X: clampInt(w/2, 1, w-2), Y: clampInt(h/2, 1, h-2)}
	m.Cells[spawn.Y][spawn.X] = engine.CellEmpty
	openCells := 1

	targetOpenCells := g.targetOpenCells(w, h)
	if targetOpenCells < minTargetOpenCells {
		targetOpenCells = minTargetOpenCells
	}

	turnChance := g.turnChance()
	x, y := spawn.X, spawn.Y
	dx, dy := randomDir(rng)
	maxSteps := w * h * maxStepsPerCell

	for steps := 0; steps < maxSteps && openCells < targetOpenCells; steps++ {
		if rng.Float64() < turnChance {
			dx, dy = randomDir(rng)
		}

		nx, ny := x+dx, y+dy
		if nx <= 0 || nx >= w-1 || ny <= 0 || ny >= h-1 {
			dx, dy = randomDir(rng)
			continue
		}

		x, y = nx, ny
		if m.Cells[y][x] == engine.CellWall {
			m.Cells[y][x] = engine.CellEmpty
			openCells++
		}
	}

	// Ensure at least two reachable tiles so stairs can be distinct from spawn.
	if openCells < 2 {
		for _, d := range []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := spawn.X+d.X, spawn.Y+d.Y
			if nx > 0 && nx < w-1 && ny > 0 && ny < h-1 {
				m.Cells[ny][nx] = engine.CellEmpty
				openCells++
				break
			}
		}
	}

	stairs := farthestReachableCell(m, spawn)
	if stairs == spawn {
		// Extremely small or unlucky maps: pick any other reachable cell.
		stairs = firstReachableDifferentCell(m, spawn)
	}
	m.Cells[stairs.Y][stairs.X] = engine.CellStairs

	g.SpawnPos = spawn
	g.StairsPos = stairs
	return m
}

func (g *FloorGenerator) targetOpenCells(w, h int) int {
	depthFactor := clamp01(float64(g.Depth) / depthScaleMax)
	openFraction := baseOpenFraction - depthFactor*openFractionDrop
	if openFraction < minOpenFraction {
		openFraction = minOpenFraction
	}
	interior := float64((w - 2) * (h - 2))
	return int(math.Round(interior * openFraction))
}

func (g *FloorGenerator) turnChance() float64 {
	depthFactor := clamp01(float64(g.Depth) / depthScaleMax)
	chance := baseTurnChance + depthFactor*turnChanceIncrease
	if chance > 1.0 {
		return 1.0
	}
	return chance
}

func newSolidWallMap(w, h int) *engine.GameMap {
	m := &engine.GameMap{
		Width:  w,
		Height: h,
		Cells:  make([][]int, h),
	}
	for y := 0; y < h; y++ {
		m.Cells[y] = make([]int, w)
		for x := 0; x < w; x++ {
			m.Cells[y][x] = engine.CellWall
		}
	}
	return m
}

func randomDir(rng *rand.Rand) (int, int) {
	switch rng.Intn(4) {
	case 0:
		return 1, 0
	case 1:
		return -1, 0
	case 2:
		return 0, 1
	default:
		return 0, -1
	}
}

func farthestReachableCell(m *engine.GameMap, from Point) Point {
	type node struct {
		p    Point
		dist int
	}

	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}

	q := make([]node, 0, m.Width*m.Height)
	head := 0
	q = append(q, node{p: from, dist: 0})
	visited[from.Y][from.X] = true

	best := from
	bestDist := 0

	for head < len(q) {
		cur := q[head]
		head++

		if cur.dist > bestDist {
			bestDist = cur.dist
			best = cur.p
		}

		for _, d := range []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := cur.p.X+d.X, cur.p.Y+d.Y
			if nx < 0 || nx >= m.Width || ny < 0 || ny >= m.Height {
				continue
			}
			if visited[ny][nx] {
				continue
			}
			if m.Cells[ny][nx] == engine.CellWall {
				continue
			}
			visited[ny][nx] = true
			q = append(q, node{p: Point{X: nx, Y: ny}, dist: cur.dist + 1})
		}
	}

	return best
}

func firstReachableDifferentCell(m *engine.GameMap, from Point) Point {
	visited := make([][]bool, m.Height)
	for y := range visited {
		visited[y] = make([]bool, m.Width)
	}

	q := make([]Point, 0, m.Width*m.Height)
	head := 0
	q = append(q, from)
	visited[from.Y][from.X] = true

	for head < len(q) {
		cur := q[head]
		head++
		if cur != from {
			return cur
		}

		for _, d := range []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			nx, ny := cur.X+d.X, cur.Y+d.Y
			if nx < 0 || nx >= m.Width || ny < 0 || ny >= m.Height {
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

	return from
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
