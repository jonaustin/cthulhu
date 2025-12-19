# CLAUDE.md

This file provides guidance for AI assistants working with code in this repository.

## Project Overview

A Lovecraftian terminal-based first-person 3D game where the player descends through procedurally-generated floors. No combat, no explicit goal — just an endless descent into cosmic horror as reality degrades around you.

## Critical Rules (TL;DR)

1. **NEVER commit to main** - Always use feature branches (`feature/<issue-id>-description`); except when closing a bd issue
2. **ALL work needs a bd issue** - `bd create --title="..." --type=task|bug|feature`
3. **ALWAYS use pull requests** - No direct pushes, no skipping review
4. **ALWAYS add tests** - No untested code

## Development Commands

### Building and Running

```bash
# Build the binary
go build -o cthulhu .

# Run with default floor size (32x32)
go run .
# or
./cthulhu

# Run with custom floor size
go run . -fs 16x16
./cthulhu -fs 48x24
```

### In-Game Controls

- **W/S** - Move forward/backward
- **A/D** - Rotate left/right
- **C** - Open cheat menu (debug/testing)
  - **T** - Teleport to specific depth
  - **+/-** - Adjust corruption bias
  - **M** - Toggle mini-map
- **Q/ESC** - Quit game

### Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run tests with coverage
go test -v -cover ./...
```

## Development Guidelines

**REQUIRED: Follow this workflow for ALL code changes.**

### Workflow Requirements

> **Note:** `bd` (beads) is a local issue tracker stored in `.beads/`. Run `bd --help` for commands.

1. **Create a bd issue** - ALL work must be tracked

   ```bash
   bd create --title="Add feature X" --type=feature
   # Returns: game-xxx
   ```

2. **Create a feature branch** - One branch per issue, NEVER commit to main

   ```bash
   # Branch naming: feature/<issue-id>-brief-description
   git checkout -b feature/game-xxx-add-feature-x

   # Update issue status
   bd update game-xxx --status=in_progress
   ```

3. **Create PR** - ALWAYS use pull requests

   ```bash
   # After committing changes
   git push -u origin feature/game-xxx-add-feature-x

   # Create PR (reference issue, but don't close yet)
   gh pr create --title "Add feature X" --body "Addresses game-xxx"
   ```

4. **Close issue AFTER merge** - Only close when PR is merged to main

   ```bash
   # After PR is merged to main
   bd close game-xxx
   bd sync                    # commits closure to main branch, pushes
   ```

**Never:**

- ❌ Write code without a bd issue; except when closing a bd issue
- ❌ Commit directly to main branch; except when closing a bd issue
- ❌ Skip pull requests; except when closing a bd issue
- ❌ Reuse branches across multiple issues
- ❌ Close a bd issue before its PR is merged to main

**Why this matters:**

- **Constants:** Consistency, one-place updates, self-documenting
- **bd tracking:** Project visibility, dependency management
- **Feature branches:** Safe parallel work, clean history
- **Pull requests:** Code review, CI checks, documented decisions

### Pre-Commit Verification

Before pushing to your feature branch:

- [ ] No magic numbers (use constants or document why)
- [ ] Imports/naming match existing patterns
- [ ] Tests added/updated for new functionality
- [ ] Error handling follows established patterns
- [ ] Branch contains only changes for this bd issue

### When Uncertain

If unsure about any of the following, **ASK before proceeding:**

- Which constant to use (or whether to create a new one)
- Whether to create a new file vs. modify existing
- How to name something (functions, variables, files)
- Whether a change warrants a new bd issue
- Architectural decisions that affect multiple packages

### If Build/Tests Fail

1. **Read the error carefully** - Understand what failed and why
2. **Check your changes** - Run `git diff` to see what you modified
3. **Fix before committing** - Never commit broken code
4. **Run tests locally** - `go test ./...` before pushing
5. **If stuck** - Describe the error clearly and ask for guidance

## Architecture

### Package Structure

See doc/ARCH.md for detailed architecture documentation.

**Current packages:**
- `main` - Game loop, event handling, HUD, cheat menu (~500 LOC)
- `engine/` - Raycaster, Player, GameMap (~600 LOC)
- `world/` - FloorGenerator, FloorManager, Corruption (~600 LOC)
- `render/` - Shading tables, corruption effects (~200 LOC)

**Total:** ~2,735 lines of Go code with comprehensive test coverage.

### Testing Strategy

ALWAYS add/update tests for new logic.

All packages have test coverage. Run tests before committing:
```bash
go test ./...           # Run all tests
go test -v ./...        # Verbose output
go test -v -cover ./... # With coverage
```

### Current Game State

**Implemented Features:**
- ✅ 3D raycasting engine (ASCII first-person view, ~60 FPS)
- ✅ Player movement (discrete WASD controls)
- ✅ Procedural floor generation (drunk walk algorithm)
- ✅ Floor transitions via discoverable stairs
- ✅ Corruption system (increases with depth ≥10)
- ✅ Visual corruption effects (glitches, whispers, fake geometry)
- ✅ HUD with depth, corruption %, controls, hints
- ✅ Mini-map overlay (toggleable)
- ✅ Cheat menu for testing

**Not Yet Implemented:**
- The Watchers (non-interactive entities)
- Smooth player movement
- Sound/audio
- Save/load system
