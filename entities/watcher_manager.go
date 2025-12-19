package entities

import (
	"math"
	"math/rand"
)

const (
	watcherSeedDepthMultiplier = int64(1_000_003)

	noiseSaltVisibility = 0xC0FFEE01
	noiseSaltGlitch     = 0xC0FFEE02
	noiseSaltGlyph      = 0xC0FFEE03
)

// WatcherManager tracks all Watchers for a floor.
type WatcherManager struct {
	Watchers []Watcher
	Depth    int
	FOV      float64
	Ticks    int
}

// NewWatcherManager creates Watchers for the given depth.
func NewWatcherManager(depth int, seed int64, fov float64) *WatcherManager {
	wm := &WatcherManager{Depth: depth, FOV: fov}
	if depth < WatcherStartDepth {
		return wm
	}
	if fov <= 0 {
		fov = defaultFOV
		wm.FOV = fov
	}

	rng := rand.New(rand.NewSource(seed + int64(depth)*watcherSeedDepthMultiplier))
	count := watcherCountForDepth(depth, rng)
	if count <= 0 {
		return wm
	}

	watchers := make([]Watcher, 0, count)
	minEdge, maxEdge := edgeOffsetRange(fov)
	for i := 0; i < count; i++ {
		watchers = append(watchers, newWatcher(rng, minEdge, maxEdge))
	}
	wm.Watchers = watchers
	return wm
}

// Update advances Watcher drift and animation ticks.
func (wm *WatcherManager) Update() {
	if wm == nil {
		return
	}
	wm.Ticks++
	if len(wm.Watchers) == 0 {
		return
	}
	fov := wm.FOV
	if fov <= 0 {
		fov = defaultFOV
	}
	minEdge, maxEdge := edgeOffsetRange(fov)
	for i := range wm.Watchers {
		w := &wm.Watchers[i]
		w.Angle += w.Drift * WatcherDriftSpeed
		if w.Angle < minEdge {
			w.Angle = minEdge
			w.Drift = 1
		}
		if w.Angle > maxEdge {
			w.Angle = maxEdge
			w.Drift = -1
		}
	}
}

// VisibleCount returns the number of Watchers visible this frame.
func (wm *WatcherManager) VisibleCount() int {
	if wm == nil {
		return 0
	}
	count := 0
	for i, w := range wm.Watchers {
		if wm.isVisible(w, i) {
			count++
		}
	}
	return count
}

// CorruptionDelta returns the corruption increment for visible Watchers this frame.
func (wm *WatcherManager) CorruptionDelta() float64 {
	return float64(wm.VisibleCount()) * WatcherCorruptionRate
}

// Sprites returns the Watcher sprites to render for the current frame.
func (wm *WatcherManager) Sprites(screenW, screenH int) []WatcherSprite {
	if wm == nil || screenW <= 0 || screenH <= 0 {
		return nil
	}
	if len(wm.Watchers) == 0 {
		return nil
	}

	fov := wm.FOV
	if fov <= 0 {
		fov = defaultFOV
	}

	sprites := make([]WatcherSprite, 0, len(wm.Watchers))
	for i, w := range wm.Watchers {
		if !wm.isVisible(w, i) {
			continue
		}
		col := watcherColumn(w, fov, screenW)
		height := watcherSpriteHeight(screenH, w.Distance)
		startY, endY := spriteSpan(screenH, height)
		if startY >= endY {
			continue
		}
		sprites = append(sprites, WatcherSprite{
			Column: col,
			StartY: startY,
			EndY:   endY,
			Char:   wm.glyphFor(w, i),
		})
	}
	return sprites
}

func watcherCountForDepth(depth int, rng *rand.Rand) int {
	if depth < WatcherStartDepth || rng == nil {
		return 0
	}
	switch {
	case depth <= watcherDepthTier1Max:
		return 1 + rng.Intn(2)
	case depth <= watcherDepthTier2Max:
		return 2 + rng.Intn(2)
	default:
		return 3 + rng.Intn(2)
	}
}

func newWatcher(rng *rand.Rand, minEdge, maxEdge float64) Watcher {
	angle := minEdge
	if maxEdge > minEdge {
		angle = minEdge + rng.Float64()*(maxEdge-minEdge)
	}
	distance := WatcherMinDistance
	if WatcherMaxDistance > WatcherMinDistance {
		distance = WatcherMinDistance + rng.Float64()*(WatcherMaxDistance-WatcherMinDistance)
	}

	side := 1
	if rng.Intn(2) == 0 {
		side = -1
	}

	drift := -1.0
	if rng.Intn(2) == 0 {
		drift = 1.0
	}

	return Watcher{
		Angle:    angle,
		Distance: distance,
		Drift:    drift,
		Side:     side,
		Seed:     uint64(rng.Int63()),
	}
}

func watcherColumn(w Watcher, fov float64, screenW int) int {
	if screenW <= 1 {
		return 0
	}
	if fov <= 0 {
		fov = defaultFOV
	}
	offset := float64(w.Side) * w.Angle
	norm := (offset / fov) + 0.5
	norm = clampFloat(norm, 0, 1)
	col := int(math.Round(norm * float64(screenW-1)))
	if col < 0 {
		return 0
	}
	if col >= screenW {
		return screenW - 1
	}
	return col
}

func (wm *WatcherManager) isVisible(w Watcher, index int) bool {
	return chance01(watcherNoise(w, wm.Ticks, index, noiseSaltVisibility), watcherVisibleChance)
}

func (wm *WatcherManager) glyphFor(w Watcher, index int) rune {
	if chance01(watcherNoise(w, wm.Ticks, index, noiseSaltGlitch), watcherGlitchChance) && len(watcherGlitchChars) > 0 {
		idx := pickIndex(watcherNoise(w, wm.Ticks, index, noiseSaltGlyph), len(watcherGlitchChars))
		return watcherGlitchChars[idx]
	}
	return watcherChar
}

func watcherNoise(w Watcher, ticks int, index int, salt uint64) uint64 {
	n := w.Seed ^ salt
	n ^= uint64(uint32(ticks)) * 0x9E3779B185EBCA87
	n ^= uint64(uint32(index)) * 0xC2B2AE3D27D4EB4F
	n ^= uint64(uint32(w.Side)) * 0x165667B19E3779F9
	return mix64(n)
}

func chance01(n uint64, p float64) bool {
	if p <= 0 {
		return false
	}
	if p >= 1 {
		return true
	}
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
