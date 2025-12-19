package render

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestApplyCharGlitchAtNoopWhenNoCorruption(t *testing.T) {
	ctx := NewEffectsContext(1, 0.0, 123)
	in := ShadeChars[0]
	if got := ApplyCharGlitchAt(in, ctx, 10, 5); got != in {
		t.Fatalf("expected no-op at 0 corruption, got %c", got)
	}
}

func TestApplyCharGlitchAtDeterministicPerFrameAndPosition(t *testing.T) {
	ctx := NewEffectsContext(20, 1.0, 777)
	in := ShadeChars[0]

	a := ApplyCharGlitchAt(in, ctx, 3, 4)
	b := ApplyCharGlitchAt(in, ctx, 3, 4)
	if a != b {
		t.Fatalf("expected deterministic output, got %c vs %c", a, b)
	}
}

func TestApplyColorBleedAtNoopWhenNoCorruption(t *testing.T) {
	ctx := NewEffectsContext(1, 0.0, 123)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	got := ApplyColorBleedAt(style, ctx, 1, 2)
	fgIn, _, _ := style.Decompose()
	fgOut, _, _ := got.Decompose()
	if fgIn != fgOut {
		t.Fatalf("expected no-op at 0 corruption, got fg %v -> %v", fgIn, fgOut)
	}
}

func TestApplyColorBleedAtDeterministicPerFrameAndPosition(t *testing.T) {
	ctx := NewEffectsContext(30, 1.0, 999)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)

	a := ApplyColorBleedAt(style, ctx, 8, 9)
	b := ApplyColorBleedAt(style, ctx, 8, 9)

	fga, _, _ := a.Decompose()
	fgb, _, _ := b.Decompose()
	if fga != fgb {
		t.Fatalf("expected deterministic output, got fg %v vs %v", fga, fgb)
	}
}

func TestVisualConfigClamp(t *testing.T) {
	ResetVisualConfig()
	defer ResetVisualConfig()

	SetVisualConfig(VisualConfig{
		VisualScale:          2.5,
		MaxCharGlitchChance:  -1.0,
		MaxColorBleedChance:  3.0,
		WhisperWindowTicks:   -10,
		MaxWhisperPerWindow:  2.0,
		MaxFakeGeometryCells: -5,
		GlitchBlinkMaxTicks:  -2,
		BleedBlinkMaxTicks:   999,
	})

	cfg := GetVisualConfig()
	if cfg.VisualScale != visualScaleMax {
		t.Fatalf("expected visual scale clamped to %f, got %f", visualScaleMax, cfg.VisualScale)
	}
	if cfg.MaxCharGlitchChance != chanceMin {
		t.Fatalf("expected char glitch clamped to %f, got %f", chanceMin, cfg.MaxCharGlitchChance)
	}
	if cfg.MaxColorBleedChance != chanceMax {
		t.Fatalf("expected color bleed clamped to %f, got %f", chanceMax, cfg.MaxColorBleedChance)
	}
	if cfg.WhisperWindowTicks != whisperWindowMin {
		t.Fatalf("expected whisper window clamped to %d, got %d", whisperWindowMin, cfg.WhisperWindowTicks)
	}
	if cfg.MaxWhisperPerWindow != chanceMax {
		t.Fatalf("expected whisper max/window clamped to %f, got %f", chanceMax, cfg.MaxWhisperPerWindow)
	}
	if cfg.MaxFakeGeometryCells != maxFakeGeoCellsMin {
		t.Fatalf("expected fake geometry cells clamped to %d, got %d", maxFakeGeoCellsMin, cfg.MaxFakeGeometryCells)
	}
	if cfg.GlitchBlinkMaxTicks != blinkMaxTicksMin {
		t.Fatalf("expected glitch blink clamped to %d, got %d", blinkMaxTicksMin, cfg.GlitchBlinkMaxTicks)
	}
	if cfg.BleedBlinkMaxTicks != blinkMaxTicksMax {
		t.Fatalf("expected bleed blink clamped to %d, got %d", blinkMaxTicksMax, cfg.BleedBlinkMaxTicks)
	}
}
