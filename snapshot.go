package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type snapshotMeta struct {
	Depth      int
	Corruption float64
	Ticks      int
	Width      int
	Height     int
	Timestamp  time.Time
}

func captureScreenLines(screen tcell.Screen, width, height int) []string {
	if screen == nil || width <= 0 || height <= 0 {
		return nil
	}

	lines := make([]string, height)
	for y := 0; y < height; y++ {
		row := make([]rune, width)
		for x := 0; x < width; x++ {
			mainc, _, _, _ := screen.GetContent(x, y)
			if mainc == 0 {
				mainc = ' '
			}
			row[x] = mainc
		}
		lines[y] = string(row)
	}
	return lines
}

func snapshotFilename(meta snapshotMeta) string {
	ts := meta.Timestamp.UTC().Format("20060102-150405.000")
	return fmt.Sprintf("depth-%02d-%s.txt", meta.Depth, ts)
}

func writeSnapshotFile(path string, meta snapshotMeta, lines []string) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# depth=%d corruption=%.4f ticks=%d size=%dx%d time=%s\n",
		meta.Depth,
		meta.Corruption,
		meta.Ticks,
		meta.Width,
		meta.Height,
		meta.Timestamp.UTC().Format(time.RFC3339Nano),
	))
	for _, line := range lines {
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func (g *Game) captureSnapshot() (string, error) {
	if g == nil || g.Screen == nil {
		return "", fmt.Errorf("screen not ready")
	}

	meta := snapshotMeta{
		Depth:      0,
		Corruption: g.Corruption,
		Ticks:      0,
		Width:      g.Width,
		Height:     g.Height,
		Timestamp:  time.Now(),
	}
	if g.Floor != nil {
		meta.Depth = g.Floor.Depth
	}
	if g.CorruptState != nil {
		meta.Ticks = g.CorruptState.Ticks
	}

	lines := captureScreenLines(g.Screen, g.Width, g.Height)
	if len(lines) == 0 {
		return "", fmt.Errorf("empty snapshot")
	}

	dir := filepath.Join(".", "snapshots")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, snapshotFilename(meta))
	if err := writeSnapshotFile(path, meta, lines); err != nil {
		return "", err
	}
	return path, nil
}
