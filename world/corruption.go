package world

const (
	corruptionStartDepth = 10
	corruptionMaxDepth   = 50
)

type Corruption struct {
	Level float64 // 0.0 to 1.0
	Bias  float64 // additive override, -1.0 to 1.0 (clamped in GetLevel)
	Depth int     // Current floor depth
	Ticks int     // Frame counter for animation
}

func NewCorruption() *Corruption {
	return &Corruption{}
}

func (c *Corruption) Update(depth int) {
	if c == nil {
		return
	}

	c.Ticks++
	c.Depth = depth
	c.Level = c.calculateLevel(depth)
}

func (c *Corruption) GetLevel() float64 {
	if c == nil {
		return 0
	}
	return clamp01(c.Level + c.Bias)
}

func (c *Corruption) GetBias() float64 {
	if c == nil {
		return 0
	}
	return clampNeg1To1(c.Bias)
}

func (c *Corruption) AdjustBias(delta float64) {
	if c == nil {
		return
	}
	c.Bias = clampNeg1To1(c.Bias + delta)
}

// calculateLevel maps depth to a corruption value.
//
// Corruption starts at floor 10 and reaches full corruption around floor 50.
func (c *Corruption) calculateLevel(depth int) float64 {
	if depth < corruptionStartDepth {
		return 0.0
	}
	span := float64(corruptionMaxDepth - corruptionStartDepth)
	if span <= 0 {
		return 1.0
	}
	return clamp01(float64(depth-corruptionStartDepth) / span)
}

func clampNeg1To1(v float64) float64 {
	if v < -1 {
		return -1
	}
	if v > 1 {
		return 1
	}
	return v
}
