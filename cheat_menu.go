package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"unicode"

	"game/render"

	"github.com/gdamore/tcell/v2"
)

type cheatMode int

const (
	cheatModeMain cheatMode = iota
	cheatModeTeleport
	cheatModeCorruptionTune
)

const (
	cheatCorruptionStep = 0.05
	cheatTeleportMaxBuf = 6
)

const (
	tuneVisualScale = iota
	tuneCharGlitchChance
	tuneColorBleedChance
	tuneWhisperWindowTicks
	tuneWhisperMaxPerWindow
	tuneFakeGeoMaxCells
	tuneCorruptionBias
	tuneParamCount
)

const (
	tuneVisualScaleStep     = 0.05
	tuneChanceStep          = 0.01
	tuneWhisperWindowStep   = 5
	tuneWhisperMaxStep      = 0.02
	tuneFakeGeoMaxCellsStep = 1
	tuneBiasStep            = 0.05
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
		if g.cheatMode == cheatModeTeleport || g.cheatMode == cheatModeCorruptionTune {
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
	case tcell.KeyUp:
		if g.cheatMode == cheatModeCorruptionTune {
			g.adjustTuneSelection(-1)
			return true
		}
	case tcell.KeyDown:
		if g.cheatMode == cheatModeCorruptionTune {
			g.adjustTuneSelection(1)
			return true
		}
	case tcell.KeyLeft:
		if g.cheatMode == cheatModeCorruptionTune {
			g.adjustTuneValue(-1)
			return true
		}
	case tcell.KeyRight:
		if g.cheatMode == cheatModeCorruptionTune {
			g.adjustTuneValue(1)
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
	case 'v', 'V':
		g.cheatMode = cheatModeCorruptionTune
		g.cheatTuneIndex = 0
		g.cheatMessage = ""
	case 'w', 'W':
		g.ShowWatchers = !g.ShowWatchers
	case 'p', 'P':
		if path, err := g.captureSnapshot(); err != nil {
			g.cheatMessage = fmt.Sprintf("Snapshot failed: %v", err)
		} else {
			g.cheatMessage = fmt.Sprintf("Snapshot saved: %s", filepath.Base(path))
		}
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
	g.cheatTuneIndex = 0
}

func (g *Game) closeCheatMenu() {
	g.cheatMenuOpen = false
	g.cheatMode = cheatModeMain
	g.cheatTeleportBuffer = g.cheatTeleportBuffer[:0]
	g.cheatMessage = ""
	g.cheatTuneIndex = 0
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
	if g.cheatMode == cheatModeCorruptionTune {
		return g.tuneMenuLines()
	}

	lines := []string{
		"CHEATS",
		fmt.Sprintf("M: Toggle map (%s)", onOff(g.ShowMiniMap)),
		"V: Tune corruption visuals",
		fmt.Sprintf("W: Toggle watchers (%s)", onOff(g.ShowWatchers)),
		"P: Snapshot frame",
		"T: Teleport to floor",
		fmt.Sprintf("+/-: Corruption bias (%.0f%%)", corruptionBiasPct),
		"C/Esc: Close",
	}
	if g.cheatMessage != "" {
		lines = append(lines, g.cheatMessage)
	}
	return lines
}

func (g *Game) tuneMenuLines() []string {
	cfg := render.GetVisualConfig()

	lines := []string{
		"CHEATS: CORRUPTION TUNING",
		"Use Up/Down to select, Left/Right to adjust",
		formatTuneLine(tuneVisualScale, g.cheatTuneIndex, fmt.Sprintf("Visual scale: %.2f", cfg.VisualScale)),
		formatTuneLine(tuneCharGlitchChance, g.cheatTuneIndex, fmt.Sprintf("Char glitch max: %.2f", cfg.MaxCharGlitchChance)),
		formatTuneLine(tuneColorBleedChance, g.cheatTuneIndex, fmt.Sprintf("Color bleed max: %.2f", cfg.MaxColorBleedChance)),
		formatTuneLine(tuneWhisperWindowTicks, g.cheatTuneIndex, fmt.Sprintf("Whisper window ticks: %d", cfg.WhisperWindowTicks)),
		formatTuneLine(tuneWhisperMaxPerWindow, g.cheatTuneIndex, fmt.Sprintf("Whisper max/window: %.2f", cfg.MaxWhisperPerWindow)),
		formatTuneLine(tuneFakeGeoMaxCells, g.cheatTuneIndex, fmt.Sprintf("Fake geo max cells: %d", cfg.MaxFakeGeometryCells)),
		formatTuneLine(tuneCorruptionBias, g.cheatTuneIndex, fmt.Sprintf("Corruption bias: %.0f%%", g.corruptionBiasPct())),
		"Esc: Back",
	}

	return lines
}

func formatTuneLine(index, selected int, text string) string {
	prefix := "  "
	if index == selected {
		prefix = "> "
	}
	return prefix + text
}

func onOff(v bool) string {
	if v {
		return "ON"
	}
	return "OFF"
}

func (g *Game) adjustTuneSelection(delta int) {
	if g == nil || delta == 0 {
		return
	}
	idx := g.cheatTuneIndex + delta
	if idx < 0 {
		idx = tuneParamCount - 1
	}
	if idx >= tuneParamCount {
		idx = 0
	}
	g.cheatTuneIndex = idx
}

func (g *Game) adjustTuneValue(delta int) {
	if g == nil || delta == 0 {
		return
	}

	cfg := render.GetVisualConfig()

	switch g.cheatTuneIndex {
	case tuneVisualScale:
		cfg.VisualScale += float64(delta) * tuneVisualScaleStep
	case tuneCharGlitchChance:
		cfg.MaxCharGlitchChance += float64(delta) * tuneChanceStep
	case tuneColorBleedChance:
		cfg.MaxColorBleedChance += float64(delta) * tuneChanceStep
	case tuneWhisperWindowTicks:
		cfg.WhisperWindowTicks += delta * tuneWhisperWindowStep
	case tuneWhisperMaxPerWindow:
		cfg.MaxWhisperPerWindow += float64(delta) * tuneWhisperMaxStep
	case tuneFakeGeoMaxCells:
		cfg.MaxFakeGeometryCells += delta * tuneFakeGeoMaxCellsStep
	case tuneCorruptionBias:
		if g.CorruptState != nil {
			g.CorruptState.AdjustBias(float64(delta) * tuneBiasStep)
		}
	}

	render.SetVisualConfig(cfg)
}

func (g *Game) corruptionBiasPct() float64 {
	if g == nil || g.CorruptState == nil {
		return 0
	}
	return g.CorruptState.GetBias() * 100
}
