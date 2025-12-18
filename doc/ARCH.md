# The Abyss — Architecture Document

## Overview
A Lovecraftian terminal-based first-person 3D game where the player descends through procedurally-generated floors. No combat, no explicit goal — just an endless descent into cosmic horror as reality degrades around you.

## Tech Stack
| Component | Choice | Rationale |
|-----------|--------|-----------|
| Language | Go | Fast enough for raycasting, simple to iterate |
| Terminal lib | tcell | Cross-platform, battle-tested, handles input/rendering |
| Resolution | 120x40 | Modern terminal size, good detail |

## Core Systems

### Raycaster Engine
Renders a first-person 3D view using classic raycasting (Wolfenstein 3D style):
- Cast one ray per screen column
- Calculate wall distance → wall height
- ASCII shading based on distance: `█ ▓ ▒ ░ .`

### Player
- Position (x, y float64)
- Direction (angle or vector)
- Field of view (~60°)
- Discrete movement: W/S forward/back, A/D rotate 90°

### Map Representation
2D grid where each cell is:
- `0` = empty space
- `1` = wall
- `2` = stairs down
- Future: doors, altars, special tiles

### Procedural Generation
Each floor generated fresh using:
- BSP, drunk walk, or cellular automata
- Depth influences parameters (deeper = weirder geometry)
- Guaranteed path from spawn to stairs

### Corruption System
- Corruption meter increases with depth
- Affects visual rendering:
  - Character substitution (walls flicker)
  - Color bleeding (ANSI glitches)
  - Fake geometry (illusory walls/doors)
  - Text whispers (fragments on screen)
- No death — corruption is purely perceptual

### The Watchers
- Entities appearing at vision edges
- Non-interactive presences
- Looking too long increases corruption

## File Structure
```
/Users/jon/code/game/
├── main.go           # Entry point, game loop
├── go.mod
├── doc/
│   └── ARCH.md       # This file
├── engine/
│   ├── raycaster.go  # Raycasting math and rendering
│   ├── player.go     # Player state and movement
│   └── map.go        # Map representation
├── world/
│   ├── generator.go  # Procedural floor generation
│   ├── floor.go      # Floor state
│   └── corruption.go # Corruption effects
├── render/
│   ├── terminal.go   # tcell abstraction
│   ├── shading.go    # ASCII shading tables
│   └── effects.go    # Visual corruption effects
└── entities/
    └── watcher.go    # The Watchers
```

## Game Loop
```
init()
  └── tcell.NewScreen()
  └── load/generate first floor

loop:
  └── handleInput()  → player movement, quit
  └── update()       → game state, corruption
  └── render()       → raycaster → screen
  └── screen.Show()
  └── sleep(16ms)    → ~60fps target

cleanup()
  └── screen.Fini()
```

## Design Decisions
- **No inventory** for MVP (add later if needed)
- **No death** — endless descent until quit
- **All Lovecraftian themes**: tentacles, cosmic void, forbidden knowledge
- **Discrete movement** to start (smooth later)
- **No sound** for MVP (hooks for later)

## MVP Milestone
1. Terminal renders 3D view via raycasting
2. Player moves through basic level
3. Stairs lead to next procedural floor
4. Depth counter increases
5. Basic corruption effects after floor ~10
