package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	minFloorSize = 5
	maxFloorSize = 256
)

func parseFloorSize(raw string) (int, int, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, 0, fmt.Errorf("empty")
	}

	parts := strings.Split(strings.ToLower(s), "x")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected WxH (e.g. 16x16), got %q", raw)
	}

	w, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid width: %w", err)
	}
	h, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("invalid height: %w", err)
	}
	if w < minFloorSize || h < minFloorSize {
		return 0, 0, fmt.Errorf("min size is %dx%d", minFloorSize, minFloorSize)
	}
	if w > maxFloorSize || h > maxFloorSize {
		return 0, 0, fmt.Errorf("max size is %dx%d", maxFloorSize, maxFloorSize)
	}
	return w, h, nil
}
