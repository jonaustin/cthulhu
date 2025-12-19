package engine

import (
	"math"

	"game/entities"
	"game/render"
	"github.com/gdamore/tcell/v2"
)

const (
	stairsMinSpriteHeight = 1
	stairsMaxSpriteHeight = 9
	DefaultFOV            = math.Pi / 3
	DefaultMaxDist        = 16.0
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
		FOV:          DefaultFOV,     // 60 degrees
		MaxDist:      DefaultMaxDist, // max view distance
	}
}

// Render draws the 3D view to the screen
func (r *Raycaster) Render(screen tcell.Screen, player *Player, gameMap *GameMap) {
	r.RenderWithEffects(screen, player, gameMap, render.EffectsContext{}, nil)
}

// RenderWithEffects draws the 3D view to the screen, applying deterministic corruption effects.
func (r *Raycaster) RenderWithEffects(screen tcell.Screen, player *Player, gameMap *GameMap, effects render.EffectsContext, watchers *entities.WatcherManager) {
	wallStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	ceilingStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkBlue)
	floorStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray)
	stairsStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow)
	watcherStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkMagenta)

	var watcherSprites []entities.WatcherSprite
	if watchers != nil {
		watcherSprites = watchers.Sprites(r.ScreenWidth, r.ScreenHeight)
	}

	// Cast a ray for each column
	for x := 0; x < r.ScreenWidth; x++ {
		// Calculate ray angle for this column
		// Map x from [0, width) to [-FOV/2, FOV/2]
		rayOffset := (float64(x)/float64(r.ScreenWidth) - 0.5) * r.FOV
		rayAngle := player.Angle + rayOffset

		wallDist, stairsDist := r.castRayWithStairs(player, gameMap, rayAngle)

		// Fix fish-eye: use perpendicular distance
		perpDist := wallDist * math.Cos(rayOffset)

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
				ch := render.ApplyCharGlitchAt(wallChar, effects, x, y)
				st := render.ApplyColorBleedAt(wallStyle, effects, x, y)
				screen.SetContent(x, y, ch, nil, st)
			} else {
				// Floor
				rowFromCenter := y - r.ScreenHeight/2
				floorChar := render.GetFloorShade(rowFromCenter, r.ScreenHeight/2)
				screen.SetContent(x, y, floorChar, nil, floorStyle)
			}
		}

		if stairsDist < wallDist && stairsDist < r.MaxDist {
			perpStairsDist := stairsDist * math.Cos(rayOffset)
			stairsHeight := stairsSpriteHeight(r.ScreenHeight, perpStairsDist)
			stairsY := r.ScreenHeight / 2
			startY := stairsY - stairsHeight/2
			endY := startY + stairsHeight
			if startY < 0 {
				startY = 0
			}
			if endY > r.ScreenHeight {
				endY = r.ScreenHeight
			}
			for y := startY; y < endY; y++ {
				st := render.ApplyColorBleedAt(stairsStyle, effects, x, y)
				screen.SetContent(x, y, render.StairsChar, nil, st)
			}
		}

		if len(watcherSprites) > 0 {
			for _, sprite := range watcherSprites {
				if sprite.Column != x {
					continue
				}
				for y := sprite.StartY; y < sprite.EndY; y++ {
					st := render.ApplyColorBleedAt(watcherStyle, effects, x, y)
					screen.SetContent(x, y, sprite.Char, nil, st)
				}
			}
		}
	}
}

// castRay uses DDA algorithm to find wall distance
func (r *Raycaster) castRay(player *Player, gameMap *GameMap, rayAngle float64) float64 {
	wallDist, _ := r.castRayWithStairs(player, gameMap, rayAngle)
	return wallDist
}

func stairsSpriteHeight(screenHeight int, perpDist float64) int {
	if perpDist <= 0 {
		return stairsMaxSpriteHeight
	}
	h := int(float64(screenHeight) / (perpDist * 2.0))
	if h < stairsMinSpriteHeight {
		return stairsMinSpriteHeight
	}
	if h > stairsMaxSpriteHeight {
		return stairsMaxSpriteHeight
	}
	return h
}

// castRayWithStairs uses DDA algorithm to find the wall distance while also tracking
// the nearest stairs tile encountered before the wall hit.
func (r *Raycaster) castRayWithStairs(player *Player, gameMap *GameMap, rayAngle float64) (wallDist, stairsDist float64) {
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
	stairsDist = math.Inf(1)

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

		if gameMap.GetCell(mapX, mapY) == CellStairs {
			dist := sideDistX - deltaDistX
			if side == 1 {
				dist = sideDistY - deltaDistY
			}
			if dist < stairsDist {
				stairsDist = dist
			}
		}

		// Check if ray hit a wall
		if gameMap.IsWall(mapX, mapY) {
			hit = true
		}

		// Safety: limit ray distance
		if sideDistX > r.MaxDist && sideDistY > r.MaxDist {
			return r.MaxDist, stairsDist
		}
	}

	// Calculate distance to wall
	if side == 0 {
		wallDist = sideDistX - deltaDistX
	} else {
		wallDist = sideDistY - deltaDistY
	}

	return wallDist, stairsDist
}

// SetScreenSize updates the screen dimensions
func (r *Raycaster) SetScreenSize(width, height int) {
	r.ScreenWidth = width
	r.ScreenHeight = height
}
