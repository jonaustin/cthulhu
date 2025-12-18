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

// DirX returns the X component of the direction vector
func (p *Player) DirX() float64 {
	return math.Cos(p.Angle)
}

// DirY returns the Y component of the direction vector
func (p *Player) DirY() float64 {
	return math.Sin(p.Angle)
}
