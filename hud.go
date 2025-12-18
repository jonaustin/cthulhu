package main

import (
	"math"

	"game/engine"
	"game/render"
)

const (
	defaultMiniMapRadius = 6
	stairsHintRadius     = 4
)

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func stairsHint(playerCellX, playerCellY, stairsX, stairsY int) string {
	dist := absInt(playerCellX-stairsX) + absInt(playerCellY-stairsY)
	if dist <= 1 {
		return "A cold draft spills from a nearby opening."
	}
	if dist <= stairsHintRadius {
		return "The air thins. Something leads down."
	}
	return ""
}

func buildMiniMap(gameMap *engine.GameMap, playerCellX, playerCellY, stairsX, stairsY, radius int) []string {
	if gameMap == nil || radius < 0 {
		return nil
	}

	size := radius*2 + 1
	lines := make([]string, 0, size)

	for dy := -radius; dy <= radius; dy++ {
		row := make([]rune, 0, size)
		for dx := -radius; dx <= radius; dx++ {
			x := playerCellX + dx
			y := playerCellY + dy

			ch := ' '
			if gameMap.IsValid(x, y) {
				cell := gameMap.GetCell(x, y)
				switch cell {
				case engine.CellWall:
					ch = '#'
				case engine.CellStairs:
					ch = render.StairsChar
				default:
					ch = '.'
				}
			}

			if x == stairsX && y == stairsY {
				ch = render.StairsChar
			}
			if x == playerCellX && y == playerCellY {
				ch = '@'
			}

			row = append(row, ch)
		}
		lines = append(lines, string(row))
	}

	return lines
}

func playerCell(p *engine.Player) (int, int) {
	if p == nil {
		return 0, 0
	}
	return int(math.Floor(p.X)), int(math.Floor(p.Y))
}
