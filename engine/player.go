package engine

import "math"

// Player represents the player's position and view direction
type Player struct {
	X, Y  float64 // position in map coordinates
	Angle float64 // view direction in radians (0 = east, π/2 = south)
}

// NewPlayer creates a player at the given position facing the given angle
func NewPlayer(x, y, angle float64) *Player {
	return &Player{X: x, Y: y, Angle: angle}
}

// NewPlayerAtCell creates a player centered in the given map cell.
func NewPlayerAtCell(cellX, cellY int, angle float64) *Player {
	return NewPlayer(float64(cellX)+0.5, float64(cellY)+0.5, angle)
}

// SetCell centers the player in the given map cell.
func (p *Player) SetCell(cellX, cellY int) {
	p.X = float64(cellX) + 0.5
	p.Y = float64(cellY) + 0.5
}

// DirX returns the X component of the direction vector
func (p *Player) DirX() float64 {
	return math.Cos(p.Angle)
}

// DirY returns the Y component of the direction vector
func (p *Player) DirY() float64 {
	return math.Sin(p.Angle)
}

const turnAngle = math.Pi / 2

// RotateLeft turns the player 90 degrees counter-clockwise.
func (p *Player) RotateLeft() {
	p.Angle = normalizeAngle(p.Angle - turnAngle)
}

// RotateRight turns the player 90 degrees clockwise.
func (p *Player) RotateRight() {
	p.Angle = normalizeAngle(p.Angle + turnAngle)
}

// MoveForward moves the player 1 cell in the direction they're facing if walkable.
func (p *Player) MoveForward(gameMap *GameMap) {
	p.step(gameMap, 1)
}

// MoveBackward moves the player 1 cell opposite the direction they're facing if walkable.
func (p *Player) MoveBackward(gameMap *GameMap) {
	p.step(gameMap, -1)
}

func (p *Player) step(gameMap *GameMap, step int) {
	if gameMap == nil {
		return
	}

	gridX := int(math.Floor(p.X))
	gridY := int(math.Floor(p.Y))

	dx, dy := cardinalStep(p.Angle)
	newX := gridX + dx*step
	newY := gridY + dy*step

	if gameMap.GetCell(newX, newY) == CellWall {
		return
	}

	p.X = float64(newX) + 0.5
	p.Y = float64(newY) + 0.5
}

func cardinalStep(angle float64) (int, int) {
	angle = normalizeAngle(angle)
	dir := int(math.Round(angle/turnAngle)) % 4
	switch dir {
	case 0: // east
		return 1, 0
	case 1: // south
		return 0, 1
	case 2: // west
		return -1, 0
	default: // north
		return 0, -1
	}
}

func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}
	// Clamp tiny negative/overflow due to float rounding.
	if angle >= 2*math.Pi {
		angle = 0
	}
	return angle
}
