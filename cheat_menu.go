package main

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/gdamore/tcell/v2"
)

type cheatMode int

const (
	cheatModeMain cheatMode = iota
	cheatModeTeleport
)

const (
	cheatCorruptionStep = 0.05
	cheatTeleportMaxBuf = 6
)

func (g *Game) handleCheatEvent(ev *tcell.EventKey) bool {
	if g == nil || ev == nil {
		return false
	}

	// Global cheat toggle.
	if ev.Key() == tcell.KeyRune {
		switch ev.Rune() {
		case 'c', 'C':
			if g.cheatMenuOpen {
				g.closeCheatMenu()
			} else {
				g.openCheatMenu()
			}
			return true
		}
	}

	// Ignore all other cheat handling when the menu is closed.
	if !g.cheatMenuOpen {
		return false
	}

	switch ev.Key() {
	case tcell.KeyEscape:
		if g.cheatMode == cheatModeTeleport {
			g.cheatMode = cheatModeMain
			g.cheatTeleportBuffer = g.cheatTeleportBuffer[:0]
			g.cheatMessage = ""
			return true
		}
		g.closeCheatMenu()
		return true
	case tcell.KeyEnter:
		if g.cheatMode == cheatModeTeleport {
			g.commitTeleport()
			return true
		}
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if g.cheatMode == cheatModeTeleport && len(g.cheatTeleportBuffer) > 0 {
			g.cheatTeleportBuffer = g.cheatTeleportBuffer[:len(g.cheatTeleportBuffer)-1]
			return true
		}
	}

	if ev.Key() != tcell.KeyRune {
		return true
	}

	r := ev.Rune()

	if g.cheatMode == cheatModeTeleport {
		if unicode.IsDigit(r) && len(g.cheatTeleportBuffer) < cheatTeleportMaxBuf {
			g.cheatTeleportBuffer = append(g.cheatTeleportBuffer, r)
			return true
		}
		return true
	}

	switch r {
	case 'q', 'Q':
		g.closeCheatMenu()
	case 'm', 'M':
		g.ShowMiniMap = !g.ShowMiniMap
	case 'w', 'W':
		g.ShowWatchers = !g.ShowWatchers
	case 't', 'T':
		g.cheatMode = cheatModeTeleport
		g.cheatTeleportBuffer = g.cheatTeleportBuffer[:0]
		g.cheatMessage = ""
	case '-', '_':
		if g.CorruptState != nil {
			g.CorruptState.AdjustBias(-cheatCorruptionStep)
		}
	case '+', '=':
		if g.CorruptState != nil {
			g.CorruptState.AdjustBias(cheatCorruptionStep)
		}
	}

	return true
}

func (g *Game) openCheatMenu() {
	g.cheatMenuOpen = true
	g.cheatMode = cheatModeMain
	g.cheatTeleportBuffer = g.cheatTeleportBuffer[:0]
	g.cheatMessage = ""
}

func (g *Game) closeCheatMenu() {
	g.cheatMenuOpen = false
	g.cheatMode = cheatModeMain
	g.cheatTeleportBuffer = g.cheatTeleportBuffer[:0]
	g.cheatMessage = ""
}

func (g *Game) commitTeleport() {
	raw := string(g.cheatTeleportBuffer)
	g.cheatTeleportBuffer = g.cheatTeleportBuffer[:0]

	depth, err := strconv.Atoi(raw)
	if err != nil || depth < 1 {
		g.cheatMessage = "Invalid depth"
		return
	}
	if !g.teleportToDepth(depth) {
		g.cheatMessage = "Teleport failed"
		return
	}

	g.cheatMessage = fmt.Sprintf("Teleported to depth %d", depth)
	g.cheatMode = cheatModeMain
}

func (g *Game) teleportToDepth(depth int) bool {
	if g == nil || g.FloorManager == nil {
		return false
	}
	if depth < 1 {
		depth = 1
	}

	f := g.FloorManager.TeleportToDepth(depth)
	if f == nil || f.Map == nil || g.Player == nil {
		return false
	}

	g.Floor = f
	g.GameMap = f.Map
	g.Player.SetCell(f.SpawnPos.X, f.SpawnPos.Y)
	return true
}

func (g *Game) renderCheatMenu() {
	if g == nil || g.Screen == nil {
		return
	}

	lines := g.cheatMenuLines()
	if len(lines) == 0 {
		return
	}

	menuW := 0
	for _, line := range lines {
		w := len([]rune(line))
		if w > menuW {
			menuW = w
		}
	}
	menuH := len(lines)
	if menuW <= 0 || menuH <= 0 {
		return
	}

	// Add some padding around the text.
	boxW := menuW + 4
	boxH := menuH + 2
	if boxW >= g.Width || boxH >= g.Height {
		return
	}

	startX := (g.Width - boxW) / 2
	startY := (g.Height - boxH) / 2

	borderStyle := tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorBlack)
	fillStyle := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	// Border.
	for x := 0; x < boxW; x++ {
		g.Screen.SetContent(startX+x, startY, '-', nil, borderStyle)
		g.Screen.SetContent(startX+x, startY+boxH-1, '-', nil, borderStyle)
	}
	for y := 0; y < boxH; y++ {
		g.Screen.SetContent(startX, startY+y, '|', nil, borderStyle)
		g.Screen.SetContent(startX+boxW-1, startY+y, '|', nil, borderStyle)
	}
	g.Screen.SetContent(startX, startY, '+', nil, borderStyle)
	g.Screen.SetContent(startX+boxW-1, startY, '+', nil, borderStyle)
	g.Screen.SetContent(startX, startY+boxH-1, '+', nil, borderStyle)
	g.Screen.SetContent(startX+boxW-1, startY+boxH-1, '+', nil, borderStyle)

	// Fill + text.
	for y := 0; y < menuH; y++ {
		// Fill line with spaces for consistent background.
		for x := 0; x < menuW; x++ {
			g.Screen.SetContent(startX+2+x, startY+1+y, ' ', nil, fillStyle)
		}
		g.drawString(startX+2, startY+1+y, lines[y], fillStyle)
	}
}

func (g *Game) cheatMenuLines() []string {
	if g == nil {
		return nil
	}

	depth := 0
	if g.Floor != nil {
		depth = g.Floor.Depth
	}

	corruptionBiasPct := 0.0
	if g.CorruptState != nil {
		corruptionBiasPct = g.CorruptState.GetBias() * 100
	}

	if g.cheatMode == cheatModeTeleport {
		return []string{
			"CHEATS: TELEPORT",
			fmt.Sprintf("Current depth: %d", depth),
			"Enter a depth and press Enter:",
			fmt.Sprintf("> %s", string(g.cheatTeleportBuffer)),
			"Esc: Back",
		}
	}

	lines := []string{
		"CHEATS",
		fmt.Sprintf("M: Toggle map (%s)", onOff(g.ShowMiniMap)),
		fmt.Sprintf("W: Toggle watchers (%s)", onOff(g.ShowWatchers)),
		"T: Teleport to floor",
		fmt.Sprintf("+/-: Corruption bias (%.0f%%)", corruptionBiasPct),
		"C/Esc: Close",
	}
	if g.cheatMessage != "" {
		lines = append(lines, g.cheatMessage)
	}
	return lines
}

func onOff(v bool) string {
	if v {
		return "ON"
	}
	return "OFF"
}
