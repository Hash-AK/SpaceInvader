package main

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Alien struct {
	X       int
	Y       int
	isAlive bool
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
	alienStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
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
	var aliens [][]Alien
	const aliensRow = 5
	const alienCol = 10
	for r := 0; r < aliensRow; r++ {
		var row []Alien
		for c := 0; c < alienCol; c++ {
			newAllien := Alien{
				X:       c * 3,
				Y:       r + 1,
				isAlive: true,
			}
			row = append(row, newAllien)

		}
		aliens = append(aliens, row)
	}
	alienDirection := 1
	alienMoveTimer := 0
	alienSpeedFactor := 1
	alienSpeed := 20
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
			alienMoveTimer = alienMoveTimer + 1
			alienSpeedFactor = alienSpeedFactor + 1
			if (alienSpeedFactor % 50) == 0 {
				alienSpeed = alienSpeed - 1
			}
			if alienMoveTimer >= alienSpeed {
				alienMoveTimer = 0
				hitEdge := false
				for r := range aliens {
					for c := range aliens[r] {
						if aliens[r][c].isAlive {
							if (aliens[r][c].X >= termWidth-1 && alienDirection == 1) || (aliens[r][c].X <= 0 && alienDirection == -1) {
								hitEdge = true
								break
							}
						}
					}
					if hitEdge {
						break
					}
				}
				if hitEdge {
					alienDirection *= -1
					for r := range aliens {
						for c := range aliens[r] {
							aliens[r][c].Y = aliens[r][c].Y + 1
						}
					}
				} else {
					for r := range aliens {
						for c := range aliens[r] {
							aliens[r][c].X = aliens[r][c].X + alienDirection
						}
					}
				}
			}
			player.Y = termHeight - 1
			playerRune = '^'
			alienRune := 'W'
			screen.Clear()
			var activeBullets []Bullet
			for i := range bullets {
				bullets[i].Y--
				if bullets[i].Y > 0 {
					screen.SetContent(bullets[i].X, bullets[i].Y, '|', nil, bulletStyle)
					activeBullets = append(activeBullets, bullets[i])
				}
			}
			for r := range aliens {
				for c := range aliens[r] {
					if aliens[r][c].isAlive {
						screen.SetContent(aliens[r][c].X, aliens[r][c].Y, alienRune, nil, alienStyle)

					}
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
