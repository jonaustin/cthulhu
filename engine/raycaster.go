package engine

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"game/render"
)

// Raycaster handles the 3D raycasting rendering
type Raycaster struct {
	ScreenWidth  int
	ScreenHeight int
	FOV          float64 // field of view in radians
	MaxDist      float64 // maximum render distance
}

// NewRaycaster creates a raycaster with the given screen dimensions
func NewRaycaster(width, height int) *Raycaster {
	return &Raycaster{
		ScreenWidth:  width,
		ScreenHeight: height,
		FOV:          math.Pi / 3, // 60 degrees
		MaxDist:      16.0,        // max view distance
	}
}

// Render draws the 3D view to the screen
func (r *Raycaster) Render(screen tcell.Screen, player *Player, gameMap *GameMap) {
	wallStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	ceilingStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkBlue)
	floorStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)

	// Cast a ray for each column
	for x := 0; x < r.ScreenWidth; x++ {
		// Calculate ray angle for this column
		// Map x from [0, width) to [-FOV/2, FOV/2]
		rayOffset := (float64(x)/float64(r.ScreenWidth) - 0.5) * r.FOV
		rayAngle := player.Angle + rayOffset

		// Cast ray and get perpendicular distance
		dist := r.castRay(player, gameMap, rayAngle)

		// Fix fish-eye: use perpendicular distance
		perpDist := dist * math.Cos(rayOffset)

		// Calculate wall height on screen
		var wallHeight int
		if perpDist > 0 {
			wallHeight = int(float64(r.ScreenHeight) / perpDist)
		} else {
			wallHeight = r.ScreenHeight
		}

		// Calculate draw start and end
		drawStart := (r.ScreenHeight - wallHeight) / 2
		drawEnd := drawStart + wallHeight

		// Clamp to screen bounds
		if drawStart < 0 {
			drawStart = 0
		}
		if drawEnd > r.ScreenHeight {
			drawEnd = r.ScreenHeight
		}

		// Get wall shading character based on distance
		wallChar := render.GetShade(perpDist, r.MaxDist)

		// Draw column
		for y := 0; y < r.ScreenHeight; y++ {
			if y < drawStart {
				// Ceiling
				screen.SetContent(x, y, render.CeilingChar, nil, ceilingStyle)
			} else if y < drawEnd {
				// Wall
				screen.SetContent(x, y, wallChar, nil, wallStyle)
			} else {
				// Floor
				rowFromCenter := y - r.ScreenHeight/2
				floorChar := render.GetFloorShade(rowFromCenter, r.ScreenHeight/2)
				screen.SetContent(x, y, floorChar, nil, floorStyle)
			}
		}
	}
}

// castRay uses DDA algorithm to find wall distance
func (r *Raycaster) castRay(player *Player, gameMap *GameMap, rayAngle float64) float64 {
	// Ray direction
	rayDirX := math.Cos(rayAngle)
	rayDirY := math.Sin(rayAngle)

	// Current map cell
	mapX := int(player.X)
	mapY := int(player.Y)

	// Length of ray from one x or y-side to next x or y-side
	deltaDistX := math.Abs(1 / rayDirX)
	deltaDistY := math.Abs(1 / rayDirY)

	// Direction to step in x or y (+1 or -1)
	var stepX, stepY int
	// Distance to next x or y gridline
	var sideDistX, sideDistY float64

	// Calculate step and initial sideDist
	if rayDirX < 0 {
		stepX = -1
		sideDistX = (player.X - float64(mapX)) * deltaDistX
	} else {
		stepX = 1
		sideDistX = (float64(mapX) + 1.0 - player.X) * deltaDistX
	}

	if rayDirY < 0 {
		stepY = -1
		sideDistY = (player.Y - float64(mapY)) * deltaDistY
	} else {
		stepY = 1
		sideDistY = (float64(mapY) + 1.0 - player.Y) * deltaDistY
	}

	// Perform DDA
	var side int // 0 for x-side, 1 for y-side
	hit := false

	for !hit {
		// Jump to next map square
		if sideDistX < sideDistY {
			sideDistX += deltaDistX
			mapX += stepX
			side = 0
		} else {
			sideDistY += deltaDistY
			mapY += stepY
			side = 1
		}

		// Check if ray hit a wall
		if gameMap.IsWall(mapX, mapY) {
			hit = true
		}

		// Safety: limit ray distance
		if sideDistX > r.MaxDist && sideDistY > r.MaxDist {
			return r.MaxDist
		}
	}

	// Calculate distance to wall
	var perpWallDist float64
	if side == 0 {
		perpWallDist = sideDistX - deltaDistX
	} else {
		perpWallDist = sideDistY - deltaDistY
	}

	return perpWallDist
}

// SetScreenSize updates the screen dimensions
func (r *Raycaster) SetScreenSize(width, height int) {
	r.ScreenWidth = width
	r.ScreenHeight = height
}
