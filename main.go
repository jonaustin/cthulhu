package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating screen: %v\n", err)
		os.Exit(1)
	}
	defer screen.Fini()

	if err := screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing screen: %v\n", err)
		os.Exit(1)
	}

	// Clear screen and show a message
	screen.Clear()
	style := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	msg := "The Abyss awaits... Press any key to quit."
	w, h := screen.Size()
	x := (w - len(msg)) / 2
	y := h / 2
	for i, r := range msg {
		screen.SetContent(x+i, y, r, nil, style)
	}
	screen.Show()

	// Wait for any key
	screen.PollEvent()
}
