package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"game/engine"
	"game/entities"
	"game/render"
	"game/world"

	"github.com/gdamore/tcell/v2"
)

// Game holds the core game state
type Game struct {
	Screen       tcell.Screen
	Running      bool
	Width        int
	Height       int
	Corruption   float64
	CorruptState *world.Corruption
	Hint         string
	events       chan tcell.Event
	Player       *engine.Player
	GameMap      *engine.GameMap
	Raycaster    *engine.Raycaster
	FloorManager *world.FloorManager
	Floor        *world.Floor

	ShowMiniMap  bool
	ShowWatchers bool

	cheatMenuOpen       bool
	cheatMode           cheatMode
	cheatTeleportBuffer []rune
	cheatMessage        string
}

func NewGame(screen tcell.Screen, floorWidth, floorHeight int) *Game {
	w, h := screen.Size()
	floorManager := world.NewFloorManagerWithSize(floorWidth, floorHeight)
	floor := floorManager.GenerateFirstFloor()
	g := &Game{
		Screen:       screen,
		Running:      true,
		Width:        w,
		Height:       h,
		Corruption:   0.0,
		CorruptState: world.NewCorruption(),
		events:       make(chan tcell.Event, 10),
		GameMap:      floor.Map,
		Raycaster:    engine.NewRaycaster(w, h),
		FloorManager: floorManager,
		Floor:        floor,
		ShowMiniMap:  true,
		ShowWatchers: true,
	}
	// Start player at floor spawn (facing north).
	g.Player = engine.NewPlayerAtCell(floor.SpawnPos.X, floor.SpawnPos.Y, -math.Pi/2)
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
		if g.handleCheatEvent(ev) {
			return
		}
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
	if g.FloorManager == nil || g.Floor == nil || g.GameMap == nil || g.Player == nil {
		return
	}

	cellX, cellY := playerCell(g.Player)

	g.Hint = stairsHint(cellX, cellY, g.Floor.StairsPos.X, g.Floor.StairsPos.Y)
	if g.GameMap.GetCell(cellX, cellY) == engine.CellStairs {
		g.Floor = g.FloorManager.DescendToNextFloor()
		g.GameMap = g.Floor.Map
		g.Player.SetCell(g.Floor.SpawnPos.X, g.Floor.SpawnPos.Y)
	}

	if g.Floor != nil && g.Floor.Watchers != nil {
		g.Floor.Watchers.Update()
	}

	depth := 0
	if g.Floor != nil {
		depth = g.Floor.Depth
	}
	if g.CorruptState != nil {
		if g.ShowWatchers && g.Floor != nil && g.Floor.Watchers != nil {
			g.CorruptState.AddExposure(g.Floor.Watchers.CorruptionDelta())
		}
		g.CorruptState.Update(depth)
		g.Corruption = g.CorruptState.GetLevel()
	}
}

// render draws the current game state to screen
func (g *Game) render() {
	g.Screen.Clear()

	// Render 3D view using raycaster
	effects := render.NewEffectsContext(0, 0, 0)
	if g.CorruptState != nil {
		effects = render.NewEffectsContext(g.CorruptState.Depth, g.CorruptState.GetLevel(), g.CorruptState.Ticks)
	}
	var watchers *entities.WatcherManager
	if g.ShowWatchers && g.Floor != nil {
		watchers = g.Floor.Watchers
	}
	g.Raycaster.RenderWithEffects(g.Screen, g.Player, g.GameMap, effects, watchers)

	// Screen-space corruption overlays (below HUD).
	render.RenderWhisperAt(g.Screen, effects, g.Width, g.Height)
	render.ApplyFakeGeometryAt(g.Screen, effects, g.Width, g.Height)

	// HUD overlay
	hudStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	playerStyle := tcell.StyleDefault.Foreground(tcell.ColorAqua).Background(tcell.ColorBlack)
	stairsStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack)
	dimStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkGray).Background(tcell.ColorBlack)

	// Status line at top
	depth := 0
	if g.Floor != nil {
		depth = g.Floor.Depth
	}
	status := fmt.Sprintf(" Depth: %d | Corruption: %.0f%% ", depth, g.Corruption*100)
	g.drawString(0, 0, status, hudStyle)

	// Mini-map (top-right, offset below status line)
	if g.ShowMiniMap && g.Floor != nil && g.GameMap != nil && g.Player != nil {
		cellX, cellY := playerCell(g.Player)
		lines := buildMiniMapRect(g.GameMap, cellX, cellY, g.Floor.StairsPos.X, g.Floor.StairsPos.Y, defaultMiniMapRadiusX, defaultMiniMapRadius)
		if len(lines) > 0 {
			mapH := len(lines)
			mapW := len([]rune(lines[0]))
			startX := miniMapStartX(g.Width, mapW)
			startY := 1
			if startY+mapH < g.Height {
				for y := 0; y < mapH; y++ {
					for x, r := range []rune(lines[y]) {
						if startX+x >= g.Width-miniMapRightMargin {
							break
						}
						style := dimStyle
						switch r {
						case '#':
							style = hudStyle
						case '@':
							style = playerStyle
						case render.StairsChar:
							style = stairsStyle
						}
						g.Screen.SetContent(startX+x, startY+y, r, nil, style)
					}
				}
			}
		}
	}

	// Controls at bottom
	controls := " W/S: Move | A/D: Turn | C: Cheats | Q: Quit "
	g.drawString(0, g.Height-1, controls, hudStyle)

	if g.Hint != "" && g.Height >= 2 {
		g.drawString(0, g.Height-2, " "+g.Hint+" ", stairsStyle)
	}

	if g.cheatMenuOpen {
		g.renderCheatMenu()
	}
}

// drawString is a helper to draw a string at x,y
func (g *Game) drawString(x, y int, str string, style tcell.Style) {
	for i, r := range str {
		g.Screen.SetContent(x+i, y, r, nil, style)
	}
}

func main() {
	fsDefault := fmt.Sprintf("%dx%d", world.DefaultMapWidth, world.DefaultMapHeight)
	floorSizeFlag := flag.String("fs", fsDefault, "floor size WxH (e.g. 16x16)")
	flag.Parse()

	floorW, floorH, err := parseFloorSize(*floorSizeFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid -fs %q: %v\n", *floorSizeFlag, err)
		os.Exit(2)
	}

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

	game := NewGame(screen, floorW, floorH)

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
