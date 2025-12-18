package render

// ShadeChars maps distance to ASCII characters for wall rendering
// Index 0 is closest (brightest), higher indices are farther (dimmer)
var ShadeChars = []rune{'█', '▓', '▒', '░', '.', ' '}

// GetShade returns the appropriate shading character for a given distance
// maxDist is the maximum render distance
func GetShade(distance, maxDist float64) rune {
	if distance <= 0 {
		return ShadeChars[0]
	}
	if distance >= maxDist {
		return ShadeChars[len(ShadeChars)-1]
	}

	// Map distance to shade index
	ratio := distance / maxDist
	index := int(ratio * float64(len(ShadeChars)-1))
	if index >= len(ShadeChars) {
		index = len(ShadeChars) - 1
	}
	return ShadeChars[index]
}

// CeilingChar is the character used for ceiling
const CeilingChar = ' '

// FloorChars are characters for floor rendering (closer = denser)
var FloorChars = []rune{'.', ':', ';', ' '}

// GetFloorShade returns floor shading based on row distance from center
func GetFloorShade(rowFromCenter, halfHeight int) rune {
	if halfHeight <= 0 {
		return FloorChars[0]
	}
	ratio := float64(rowFromCenter) / float64(halfHeight)
	index := int(ratio * float64(len(FloorChars)-1))
	if index < 0 {
		index = 0
	}
	if index >= len(FloorChars) {
		index = len(FloorChars) - 1
	}
	return FloorChars[index]
}
