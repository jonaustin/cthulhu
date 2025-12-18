package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Game holds the core game state
type Game struct {
	Screen       tcell.Screen
	Running      bool
	Width        int
	Height       int
	CurrentFloor int
	Corruption   float64
	events       chan tcell.Event
	// Player and Map will be added when those systems are implemented
}

func NewGame(screen tcell.Screen) *Game {
	w, h := screen.Size()
	g := &Game{
		Screen:       screen,
		Running:      true,
		Width:        w,
		Height:       h,
		CurrentFloor: 1,
		Corruption:   0.0,
		events:       make(chan tcell.Event, 10),
	}
	// Start event polling goroutine
	go g.pollEvents()
	return g
}

// pollEvents runs in a goroutine to feed events to the channel
func (g *Game) pollEvents() {
	for {
		ev := g.Screen.PollEvent()
		if ev == nil {
			return // Screen was finalized
		}
		g.events <- ev
	}
}

// handleInput processes input events non-blocking
func (g *Game) handleInput() {
	// Drain all available events
	for {
		select {
		case ev := <-g.events:
			g.processEvent(ev)
		default:
			return
		}
	}
}

func (g *Game) processEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			g.Running = false
		case tcell.KeyRune:
			switch ev.Rune() {
			case 'q', 'Q':
				g.Running = false
			case 'w', 'W':
				// Move forward (placeholder)
			case 's', 'S':
				// Move backward (placeholder)
			case 'a', 'A':
				// Rotate left (placeholder)
			case 'd', 'D':
				// Rotate right (placeholder)
			}
		}
	case *tcell.EventResize:
		g.Width, g.Height = ev.Size()
		g.Screen.Sync()
	}
}

// update processes game state changes
func (g *Game) update() {
	// Placeholder for game logic:
	// - Player movement
	// - Corruption increase
	// - Floor transitions
	// - Watcher spawning
}

// render draws the current game state to screen
func (g *Game) render() {
	g.Screen.Clear()

	// Placeholder: render basic info until raycaster is implemented
	style := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	dimStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGreen)

	// Title
	title := "THE ABYSS"
	g.drawString(g.Width/2-len(title)/2, 1, title, style)

	// Status line
	status := fmt.Sprintf("Floor: %d | Corruption: %.1f%%", g.CurrentFloor, g.Corruption*100)
	g.drawString(g.Width/2-len(status)/2, 3, status, dimStyle)

	// Placeholder for raycaster view
	viewMsg := "[ Raycaster view will render here ]"
	g.drawString(g.Width/2-len(viewMsg)/2, g.Height/2, viewMsg, dimStyle)

	// Controls
	controls := "W/S: Move | A/D: Turn | Q: Quit"
	g.drawString(g.Width/2-len(controls)/2, g.Height-2, controls, dimStyle)
}

// drawString is a helper to draw a string at x,y
func (g *Game) drawString(x, y int, str string, style tcell.Style) {
	for i, r := range str {
		g.Screen.SetContent(x+i, y, r, nil, style)
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating screen: %v\n", err)
		os.Exit(1)
	}

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	game := NewGame(screen)

	// Main game loop - target ~60fps (16ms per frame)
	frameDuration := time.Duration(16) * time.Millisecond

	for game.Running {
		frameStart := time.Now()

		game.handleInput()
		game.update()
		game.render()
		game.Screen.Show()

		// Sleep for remaining frame time
		elapsed := time.Since(frameStart)
		if elapsed < frameDuration {
			time.Sleep(frameDuration - elapsed)
		}
	}
}
