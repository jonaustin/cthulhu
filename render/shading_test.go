package render

import "testing"

func TestGetShade(t *testing.T) {
	maxDist := 16.0

	// Very close should be solid block
	if GetShade(0.1, maxDist) != '█' {
		t.Error("Very close distance should return solid block")
	}

	// Very far should be space
	if GetShade(maxDist, maxDist) != ' ' {
		t.Error("Max distance should return space")
	}

	// Distance 0 or negative should return solid
	if GetShade(0, maxDist) != '█' {
		t.Error("Zero distance should return solid block")
	}
	if GetShade(-1, maxDist) != '█' {
		t.Error("Negative distance should return solid block")
	}

	// Beyond max should return space
	if GetShade(maxDist+10, maxDist) != ' ' {
		t.Error("Beyond max distance should return space")
	}
}

func TestGetFloorShade(t *testing.T) {
	halfHeight := 20

	// At horizon (row 0) should be densest
	shade := GetFloorShade(0, halfHeight)
	if shade != FloorChars[0] {
		t.Errorf("At horizon expected %c, got %c", FloorChars[0], shade)
	}

	// At bottom should be sparse
	shade = GetFloorShade(halfHeight, halfHeight)
	if shade != FloorChars[len(FloorChars)-1] {
		t.Errorf("At bottom expected %c, got %c", FloorChars[len(FloorChars)-1], shade)
	}
}

func TestShadeCharsOrder(t *testing.T) {
	// Verify the shade characters are in expected order (dense to sparse)
	expected := []rune{'█', '▓', '▒', '░', '.', ' '}
	if len(ShadeChars) != len(expected) {
		t.Errorf("Expected %d shade chars, got %d", len(expected), len(ShadeChars))
	}
	for i, c := range expected {
		if ShadeChars[i] != c {
			t.Errorf("ShadeChars[%d] expected %c, got %c", i, c, ShadeChars[i])
		}
	}
}
