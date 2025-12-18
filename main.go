package main

import (
	"fmt"
	"os"
	"time"

	"game/engine"

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
	Player       *engine.Player
	GameMap      *engine.GameMap
	Raycaster    *engine.Raycaster
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
		GameMap:      engine.NewTestMap(),
		Raycaster:    engine.NewRaycaster(w, h),
	}
	// Start player in open area of test map (position 8, 8 facing north)
	g.Player = engine.NewPlayer(8.5, 8.5, -3.14159/2) // facing north (up)
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
				g.Player.MoveForward(g.GameMap)
			case 's', 'S':
				g.Player.MoveBackward(g.GameMap)
			case 'a', 'A':
				g.Player.RotateLeft()
			case 'd', 'D':
				g.Player.RotateRight()
			}
		}
	case *tcell.EventResize:
		g.Width, g.Height = ev.Size()
		g.Raycaster.SetScreenSize(g.Width, g.Height)
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

	// Render 3D view using raycaster
	g.Raycaster.Render(g.Screen, g.Player, g.GameMap)

	// HUD overlay
	hudStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)

	// Status line at top
	status := fmt.Sprintf(" Floor: %d | Corruption: %.0f%% ", g.CurrentFloor, g.Corruption*100)
	g.drawString(0, 0, status, hudStyle)

	// Controls at bottom
	controls := " W/S: Move | A/D: Turn | Q: Quit "
	g.drawString(0, g.Height-1, controls, hudStyle)
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
