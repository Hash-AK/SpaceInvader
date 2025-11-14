package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

type Alien struct {
	X     int
	Y     int
	Speed int
	Size  int
}
type Player struct {
	X int
	Y int
}

var playerRune rune

func main() {
	screen, err := tcell.NewScreen()

	if err != nil {
		fmt.Printf("Error initializing screen : %s", err)
	}
	defer screen.Fini()
	screen.Init()
	termWidth, termHeight := screen.Size()

	player := Player{
		X: termWidth / 2,
		Y: termHeight - 1,
	}
	for {
		player.Y = termHeight - 1
		playerRune = '^'
		screen.Clear()

		screen.SetContent(player.X, player.Y, playerRune, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
		screen.Show()
		event := screen.PollEvent()
		switch ev := event.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			}
			if ev.Key() == tcell.KeyLeft {
				if player.X > 0 {
					player.X = player.X - 1
				}
			}
			if ev.Key() == tcell.KeyRight {
				if player.X < termWidth-1 {
					player.X = player.X + 1
				}
			}
		case *tcell.EventResize:
			screen.Clear()
			termWidth, termHeight = screen.Size()

		}
	}
}
