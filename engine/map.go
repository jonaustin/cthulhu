package engine

// Cell type constants
const (
	CellEmpty  = 0 // walkable space
	CellWall   = 1 // solid wall
	CellStairs = 2 // stairs down (for later)
)

// GameMap represents a 2D grid-based level
type GameMap struct {
	Width  int
	Height int
	Cells  [][]int
}

// NewTestMap creates a hardcoded 16x16 test map for raycaster development
func NewTestMap() *GameMap {
	m := &GameMap{
		Width:  16,
		Height: 16,
		Cells:  make([][]int, 16),
	}

	// Initialize all cells as empty
	for y := 0; y < 16; y++ {
		m.Cells[y] = make([]int, 16)
	}

	// Test map layout:
	// ################
	// #..............#
	// #..............#
	// #..###.........#
	// #..#...........#
	// #..#...........#
	// #..............#
	// #.......####...#
	// #.......#......#
	// #.......#......#
	// #..............#
	// #.........##...#
	// #..........#...#
	// #..............#
	// #..............#
	// ################

	// Define the map using a string grid for clarity
	layout := []string{
		"################",
		"#..............#",
		"#..............#",
		"#..###.........#",
		"#..#...........#",
		"#..#...........#",
		"#..............#",
		"#.......####...#",
		"#.......#......#",
		"#.......#......#",
		"#..............#",
		"#.........##...#",
		"#..........#...#",
		"#..............#",
		"#..............#",
		"################",
	}

	for y, row := range layout {
		for x, c := range row {
			if c == '#' {
				m.Cells[y][x] = CellWall
			} else {
				m.Cells[y][x] = CellEmpty
			}
		}
	}

	return m
}

// IsValid checks if coordinates are within map bounds
func (m *GameMap) IsValid(x, y int) bool {
	return x >= 0 && x < m.Width && y >= 0 && y < m.Height
}

// GetCell returns the cell value at (x, y), or CellWall if out of bounds
func (m *GameMap) GetCell(x, y int) int {
	if !m.IsValid(x, y) {
		return CellWall // treat out-of-bounds as wall
	}
	return m.Cells[y][x]
}

// IsWall returns true if the cell at (x, y) is a wall
func (m *GameMap) IsWall(x, y int) bool {
	return m.GetCell(x, y) == CellWall
}
