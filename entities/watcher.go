package entities

import "math"

const (
	// WatcherStartDepth is the floor depth where Watchers begin appearing.
	WatcherStartDepth = 15

	// WatcherMinDistance/WatcherMaxDistance bound the perceived distance of Watchers.
	WatcherMinDistance = 8.0
	WatcherMaxDistance = 14.0

	// WatcherEdgeThreshold is the fraction of the screen reserved for edge sightings.
	WatcherEdgeThreshold = 0.20

	// WatcherDriftSpeed controls how quickly Watchers slide along the vision edge.
	WatcherDriftSpeed = 0.002

	// WatcherCorruptionRate is the per-frame corruption gain per visible Watcher.
	WatcherCorruptionRate = 0.0001
)

const (
	watcherDepthTier1Max = 24
	watcherDepthTier2Max = 34

	watcherMinSpriteHeight = 2
	watcherMaxSpriteHeight = 10

	watcherVisibleChance = 0.70
	watcherGlitchChance  = 0.30

	watcherChar = 'W'
)

var watcherGlitchChars = []rune{'#', '%', '&', '@', 'X'}

const (
	defaultFOV = math.Pi / 3
)

// Watcher represents a single edge-of-vision entity.
type Watcher struct {
	Angle    float64
	Distance float64
	Drift    float64
	Side     int
	Seed     uint64
}

// WatcherSprite describes how to draw a Watcher for the current frame.
type WatcherSprite struct {
	Column int
	StartY int
	EndY   int
	Char   rune
}

func edgeOffsetRange(fov float64) (float64, float64) {
	if fov <= 0 {
		fov = defaultFOV
	}
	half := fov * 0.5
	min := half * (1.0 - 2.0*WatcherEdgeThreshold)
	if min < 0 {
		min = 0
	}
	return min, half
}

func watcherSpriteHeight(screenHeight int, distance float64) int {
	if screenHeight <= 0 {
		return 0
	}
	if distance <= 0 {
		return watcherMaxSpriteHeight
	}
	h := int(float64(screenHeight) / (distance * 2.0))
	if h < watcherMinSpriteHeight {
		return watcherMinSpriteHeight
	}
	if h > watcherMaxSpriteHeight {
		return watcherMaxSpriteHeight
	}
	return h
}

func spriteSpan(screenHeight, spriteHeight int) (int, int) {
	if screenHeight <= 0 || spriteHeight <= 0 {
		return 0, 0
	}
	center := screenHeight / 2
	start := center - spriteHeight/2
	end := start + spriteHeight
	if start < 0 {
		start = 0
	}
	if end > screenHeight {
		end = screenHeight
	}
	return start, end
}

func clampFloat(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
