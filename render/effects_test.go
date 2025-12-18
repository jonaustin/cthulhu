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
