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
├── main.go           # Entry point, game loop, event handling
├── cheat_menu.go     # Debug/testing cheat menu (C key)
├── hud.go            # HUD rendering, mini-map, stairs hints
├── flags.go          # CLI flag parsing (floor size)
├── go.mod
├── doc/
│   └── ARCH.md       # This file
├── engine/
│   ├── raycaster.go  # Raycasting math and 3D rendering
│   ├── player.go     # Player state and movement
│   └── map.go        # Map representation (2D grid)
├── entities/
│   ├── watcher.go    # Watcher entity definitions
│   └── watcher_manager.go # Watcher spawning + drift
├── world/
│   ├── generator.go  # Procedural floor generation (drunk walk)
│   ├── floor.go      # Floor state and FloorManager
│   └── corruption.go # Corruption level calculation
└── render/
    ├── shading.go    # ASCII shading tables (walls, floors)
    └── effects.go    # Visual corruption effects (glitch, whispers, fake geo)
```

**Note:** The `entities/` package houses The Watchers (edge-of-vision entities).

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

## Implementation Status

### ✅ Completed (MVP Achieved)
1. ✅ Terminal renders 3D view via raycasting (60 FPS target)
2. ✅ Player moves through procedurally-generated levels (WASD controls)
3. ✅ Stairs lead to next floor with discoverable hints
4. ✅ Depth counter and corruption tracking
5. ✅ Corruption effects starting at depth 10:
   - Character glitching (walls flicker)
   - Color bleeding (ANSI color shifts)
   - Whispers (text fragments at 65%+ corruption)
   - Fake geometry (illusory walls at 90%+ corruption)
6. ✅ HUD with depth, corruption %, controls, stairs hints
7. ✅ Mini-map overlay (toggleable via cheat menu)
8. ✅ Configurable floor size (`-fs WxH` flag)
9. ✅ The Watchers (edge-of-vision presences)
10. ✅ Comprehensive test coverage

### 🚧 Future Enhancements
- Smooth player movement (currently discrete)
- Sound/audio hooks
- Additional corruption effects
- Save/load system
