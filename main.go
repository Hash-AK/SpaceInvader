package main

import (
	"fmt"
	"time"

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
type Bullet struct {
	X        int
	Y        int
	isActive bool
}

var playerRune rune

func main() {
	bulletStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)

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
	var bullets []Bullet
	eventChan := make(chan tcell.Event)
	quitChan := make(chan struct{})
	go func() {
		for {
			event := screen.PollEvent()
			eventChan <- event
		}
	}()
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()
Loop:
	for {
		select {
		case event := <-eventChan:
			switch ev := event.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
					close(quitChan)
				}
				if ev.Key() == tcell.KeyLeft && player.X > 0 {
					player.X = player.X - 1
				}
				if ev.Key() == tcell.KeyRight && player.X < termWidth-1 {
					player.X = player.X + 1
				}
				if ev.Key() == tcell.KeyRune {
					if ev.Rune() == ' ' {
						newBullet := Bullet{
							X:        player.X,
							Y:        player.Y - 1,
							isActive: true,
						}
						bullets = append(bullets, newBullet)
					}
				}
			case *tcell.EventResize:
				screen.Sync()
				termWidth, termHeight = screen.Size()
			}
		case <-ticker.C:
			player.Y = termHeight - 1
			playerRune = '^'
			screen.Clear()
			var activeBullets []Bullet
			for i := range bullets {
				bullets[i].Y--
				if bullets[i].Y > 0 {
					screen.SetContent(bullets[i].X, bullets[i].Y, '|', nil, bulletStyle)
					activeBullets = append(activeBullets, bullets[i])
				}
			}
			bullets = activeBullets
			screen.SetContent(player.X, player.Y, playerRune, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
			screen.Show()

		case <-quitChan:
			break Loop
		}

	}
}
