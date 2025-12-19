package render

const (
	defaultMaxCharGlitchChance   = 0.10
	defaultMaxColorBleedChance   = 0.05
	defaultWhisperWindowTicks    = 45
	defaultMaxWhisperPerWindow   = 0.12
	defaultMaxFakeGeometryCells  = 24
	defaultCorruptionVisualScale = 0.5
	defaultGlitchBlinkMaxTicks   = 24
	defaultBleedBlinkMaxTicks    = 18
)

const (
	visualScaleMin = 0.0
	visualScaleMax = 1.0

	chanceMin = 0.0
	chanceMax = 1.0

	whisperWindowMin = 1
	whisperWindowMax = 300

	maxFakeGeoCellsMin = 0
	maxFakeGeoCellsMax = 200

	blinkMaxTicksMin = 1
	blinkMaxTicksMax = 120
)

// VisualConfig controls corruption visuals at runtime.
type VisualConfig struct {
	VisualScale          float64
	MaxCharGlitchChance  float64
	MaxColorBleedChance  float64
	WhisperWindowTicks   int
	MaxWhisperPerWindow  float64
	MaxFakeGeometryCells int
	GlitchBlinkMaxTicks  int
	BleedBlinkMaxTicks   int
}

var defaultVisualConfig = VisualConfig{
	VisualScale:          defaultCorruptionVisualScale,
	MaxCharGlitchChance:  defaultMaxCharGlitchChance,
	MaxColorBleedChance:  defaultMaxColorBleedChance,
	WhisperWindowTicks:   defaultWhisperWindowTicks,
	MaxWhisperPerWindow:  defaultMaxWhisperPerWindow,
	MaxFakeGeometryCells: defaultMaxFakeGeometryCells,
	GlitchBlinkMaxTicks:  defaultGlitchBlinkMaxTicks,
	BleedBlinkMaxTicks:   defaultBleedBlinkMaxTicks,
}

var visualConfig = defaultVisualConfig

// GetVisualConfig returns the current corruption visual tuning values.
func GetVisualConfig() VisualConfig {
	return visualConfig
}

// SetVisualConfig updates the current corruption visual tuning values.
func SetVisualConfig(cfg VisualConfig) {
	visualConfig = sanitizeVisualConfig(cfg)
}

// ResetVisualConfig restores default corruption visual tuning values.
func ResetVisualConfig() {
	visualConfig = defaultVisualConfig
}

func sanitizeVisualConfig(cfg VisualConfig) VisualConfig {
	cfg.VisualScale = clampFloat(cfg.VisualScale, visualScaleMin, visualScaleMax)
	cfg.MaxCharGlitchChance = clampFloat(cfg.MaxCharGlitchChance, chanceMin, chanceMax)
	cfg.MaxColorBleedChance = clampFloat(cfg.MaxColorBleedChance, chanceMin, chanceMax)
	cfg.MaxWhisperPerWindow = clampFloat(cfg.MaxWhisperPerWindow, chanceMin, chanceMax)
	cfg.WhisperWindowTicks = clampInt(cfg.WhisperWindowTicks, whisperWindowMin, whisperWindowMax)
	cfg.MaxFakeGeometryCells = clampInt(cfg.MaxFakeGeometryCells, maxFakeGeoCellsMin, maxFakeGeoCellsMax)
	cfg.GlitchBlinkMaxTicks = clampInt(cfg.GlitchBlinkMaxTicks, blinkMaxTicksMin, blinkMaxTicksMax)
	cfg.BleedBlinkMaxTicks = clampInt(cfg.BleedBlinkMaxTicks, blinkMaxTicksMin, blinkMaxTicksMax)
	return cfg
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

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
