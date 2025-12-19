package render

import (
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Randomly replace wall chars with glitchy alternatives.
var glitchChars = []rune{'╳', '◊', '∆', '¤', '§', '░'}

// ANSI color shifts - walls occasionally wrong color.
var corruptColors = []tcell.Color{tcell.ColorRed, tcell.ColorFuchsia, tcell.ColorDarkRed}

// Lovecraftian text fragments appearing on screen.
var whispers = []string{
	"ph'nglui mglw'nafh",
	"the deep calls",
	"THEY SEE YOU",
	"do not look back",
	"the angles are wrong",
}

const (
	whisperStartLevel = 0.65
	fakeGeoStartLevel = 0.90
)

type EffectsContext struct {
	Corruption float64
	Depth      int
	Ticks      int
	Seed       uint64
}

func NewEffectsContext(depth int, corruption float64, ticks int) EffectsContext {
	return EffectsContext{
		Corruption: corruption,
		Depth:      depth,
		Ticks:      ticks,
		Seed:       mix64(uint64(depth)),
	}
}

// ApplyCharGlitch randomly replaces wall chars with glitchy alternatives.
//
// This helper is intentionally context-free; for deterministic (per-frame + position)
// effects, prefer ApplyCharGlitchAt.
func ApplyCharGlitch(char rune, corruption float64) rune {
	ctx := EffectsContext{
		Corruption: corruption,
		Seed:       uint64(time.Now().UnixNano()),
		Ticks:      int(time.Now().UnixNano()),
	}
	return ApplyCharGlitchAt(char, ctx, 0, 0)
}

func ApplyCharGlitchAt(char rune, ctx EffectsContext, x, y int) rune {
	if ctx.Corruption <= 0 || !isWallShadeChar(char) {
		return char
	}

	cfg := GetVisualConfig()
	ctxBlink := ctx
	ctxBlink.Ticks = ctx.Ticks / blinkTicksFor(ctx.Corruption, cfg.GlitchBlinkMaxTicks)
	p := clamp01(ctx.Corruption) * cfg.MaxCharGlitchChance * cfg.VisualScale
	if !chance01(cellNoise(ctxBlink, x, y, 0xA11CE), p) {
		return char
	}

	return glitchChars[pickIndex(cellNoise(ctxBlink, x, y, 0xC0FFEE), len(glitchChars))]
}

// ApplyColorBleed occasionally shifts the wall color.
//
// This helper is intentionally context-free; for deterministic (per-frame + position)
// effects, prefer ApplyColorBleedAt.
func ApplyColorBleed(style tcell.Style, corruption float64) tcell.Style {
	ctx := EffectsContext{
		Corruption: corruption,
		Seed:       uint64(time.Now().UnixNano()),
		Ticks:      int(time.Now().UnixNano()),
	}
	return ApplyColorBleedAt(style, ctx, 0, 0)
}

func ApplyColorBleedAt(style tcell.Style, ctx EffectsContext, x, y int) tcell.Style {
	if ctx.Corruption <= 0 {
		return style
	}

	cfg := GetVisualConfig()
	ctxBlink := ctx
	ctxBlink.Ticks = ctx.Ticks / blinkTicksFor(ctx.Corruption, cfg.BleedBlinkMaxTicks)
	p := clamp01(ctx.Corruption) * cfg.MaxColorBleedChance * cfg.VisualScale
	if !chance01(cellNoise(ctxBlink, x, y, 0xB1EED), p) {
		return style
	}

	return style.Foreground(corruptColors[pickIndex(cellNoise(ctxBlink, x, y, 0xD15EA5E), len(corruptColors))])
}

func RenderWhisper(screen tcell.Screen, corruption float64) {
	if screen == nil {
		return
	}
	w, h := screen.Size()
	RenderWhisperAt(screen, EffectsContext{
		Corruption: corruption,
		Seed:       uint64(time.Now().UnixNano()),
		Ticks:      int(time.Now().UnixNano()),
	}, w, h)
}

func RenderWhisperAt(screen tcell.Screen, ctx EffectsContext, width, height int) {
	if screen == nil || width <= 0 || height <= 0 {
		return
	}
	if ctx.Corruption < whisperStartLevel || len(whispers) == 0 {
		return
	}

	cfg := GetVisualConfig()
	intensity := clamp01((ctx.Corruption - whisperStartLevel) / (1.0 - whisperStartLevel))
	intensity *= cfg.VisualScale
	window := 0
	if cfg.WhisperWindowTicks > 0 {
		window = ctx.Ticks / cfg.WhisperWindowTicks
	}

	if !chance01(mix64(ctx.Seed^uint64(window)^0x51A57E), intensity*cfg.MaxWhisperPerWindow) {
		return
	}

	msg := whispers[pickIndex(mix64(ctx.Seed^uint64(window)^0x57EAD), len(whispers))]
	if len(msg) == 0 {
		return
	}

	usableHeight := height - 3
	if usableHeight <= 0 {
		return
	}
	y := 1 + pickIndex(mix64(ctx.Seed^uint64(window)^0x900D), usableHeight)
	maxX := width - len([]rune(msg))
	if maxX < 0 {
		maxX = 0
	}
	x := pickIndex(mix64(ctx.Seed^uint64(window)^0xBADC0DE), maxX+1)

	style := tcell.StyleDefault.Foreground(tcell.ColorDarkMagenta).Background(tcell.ColorBlack)
	for i, r := range []rune(msg) {
		if x+i >= width {
			break
		}
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func ApplyFakeGeometry(screen tcell.Screen, corruption float64) {
	if screen == nil {
		return
	}
	w, h := screen.Size()
	ApplyFakeGeometryAt(screen, EffectsContext{
		Corruption: corruption,
		Seed:       uint64(time.Now().UnixNano()),
		Ticks:      int(time.Now().UnixNano()),
	}, w, h)
}

func ApplyFakeGeometryAt(screen tcell.Screen, ctx EffectsContext, width, height int) {
	if screen == nil || width <= 0 || height <= 0 {
		return
	}
	if ctx.Corruption < fakeGeoStartLevel {
		return
	}

	cfg := GetVisualConfig()
	intensity := clamp01((ctx.Corruption - fakeGeoStartLevel) / (1.0 - fakeGeoStartLevel))
	intensity *= cfg.VisualScale
	count := int(intensity * float64(cfg.MaxFakeGeometryCells))
	if count < 1 {
		count = 1
	}

	style := tcell.StyleDefault.Foreground(tcell.ColorDarkRed).Background(tcell.ColorBlack)
	char := '▒'

	for i := 0; i < count; i++ {
		n := mix64(ctx.Seed ^ uint64(ctx.Ticks) ^ uint64(i)*0x9E3779B97F4A7C15)
		x := pickIndex(n^0x1234, width)
		y := pickIndex(n^0xBEEF, height)

		// Avoid painting over HUD areas by staying away from the very top/bottom rows.
		if y <= 1 || y >= height-2 {
			continue
		}

		screen.SetContent(x, y, char, nil, style)
	}
}

func isWallShadeChar(r rune) bool {
	for _, c := range ShadeChars {
		if r == c {
			return true
		}
	}
	return false
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

func blinkTicksFor(corruption float64, maxTicks int) int {
	if maxTicks <= 1 {
		return 1
	}
	c := clamp01(corruption)
	if c <= 0 {
		return maxTicks
	}
	if c >= 1 {
		return 1
	}
	return int(math.Round(float64(maxTicks) - c*float64(maxTicks-1)))
}

func cellNoise(ctx EffectsContext, x, y int, salt uint64) uint64 {
	n := ctx.Seed ^ salt
	n ^= uint64(uint32(ctx.Ticks)) * 0xD2B74407B1CE6E93
	n ^= uint64(uint32(ctx.Depth)) * 0x165667B19E3779F9
	n ^= uint64(uint32(x)) * 0x9E3779B185EBCA87
	n ^= uint64(uint32(y)) * 0xC2B2AE3D27D4EB4F
	return mix64(n)
}

func chance01(n uint64, p float64) bool {
	if p <= 0 {
		return false
	}
	if p >= 1 {
		return true
	}
	// Convert to [0,1) using the top 53 bits.
	f := float64(n>>11) * (1.0 / (1 << 53))
	return f < p
}

func pickIndex(n uint64, size int) int {
	if size <= 0 {
		return 0
	}
	return int(n % uint64(size))
}

// mix64 is a SplitMix64-style mixer.
func mix64(x uint64) uint64 {
	x += 0x9E3779B97F4A7C15
	x = (x ^ (x >> 30)) * 0xBF58476D1CE4E5B9
	x = (x ^ (x >> 27)) * 0x94D049BB133111EB
	return x ^ (x >> 31)
}
