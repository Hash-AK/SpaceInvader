package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Asteroid struct {
	X           int
	Y           int
	Orientation int
	Speed       int
}
type Player struct {
	X           int
	Y           int
	Orientation int
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Printf("Error initializing screen : %s", err)
	}
	defer screen.Fini()
	screen.Clear()
}
