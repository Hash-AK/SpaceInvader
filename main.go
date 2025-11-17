package main

import (
	"fmt"
	"math/rand"
	"os"
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

func printString(screen tcell.Screen, x int, y int, style tcell.Style, toPrint string) {
	for _, r := range toPrint {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

var score int

func checkAlienDeaths(aliens [][]Alien) (bool, int) {
	totalAliens := 0
	var cols int
	won := false
	var totalAliveAliens int

	for r := range aliens {
		for c := range aliens[r] {
			cols = aliens[r][c].X
			cols++
			totalAliens = totalAliens + 1

		}
	}
	if totalAliens > 0 {
		totalAliveAliens = totalAliens
		for r := range aliens {
			for c := range aliens[r] {
				if !aliens[r][c].isAlive {
					totalAliveAliens = totalAliveAliens - 1
				}
			}
		}
		if totalAliveAliens == 0 {
			won = true
		}
	}
	return won, totalAliveAliens

}
func main() {
	score = 0
	var gameState string
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
	lives := 3
	gameState = "playing"
	var bullets []Bullet
	var aliens [][]Alien
	var alienBullets []Bullet

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
	rand.New(rand.NewSource(time.Now().UnixNano()))
	const alienShootProb = 5

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
			if gameState == "playing" {
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
			} else {
				switch ev := event.(type) {
				case *tcell.EventKey:
					if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
						close(quitChan)
					}
				}
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
			if rand.Intn(100) < alienShootProb {
				var shooters []Alien
				for c := 0; c < alienCol; c++ {
					for r := aliensRow - 1; r >= 0; r-- {
						if aliens[r][c].isAlive {
							shooters = append(shooters, aliens[r][c])
							break
						}
					}
				}
				if len(shooters) > 0 {
					shooter := shooters[rand.Intn(len(shooters))]
					newBullet := Bullet{
						X:        shooter.X,
						Y:        shooter.Y + 1,
						isActive: true,
					}
					alienBullets = append(alienBullets, newBullet)
				}
			}
			player.Y = termHeight - 1
			playerRune = '^'
			alienRune := 'W'
			screen.Clear()
			var activeBullets []Bullet
			var activeAlienBullets []Bullet
			for i := range bullets {
				bullets[i].Y--
				if bullets[i].Y > 0 {
					bullets[i].isActive = true

					// Check collision first
					hitAlien := false
					for r := range aliens {
						for c := range aliens[r] {
							if bullets[i].X == aliens[r][c].X && bullets[i].Y == aliens[r][c].Y && aliens[r][c].isAlive {
								aliens[r][c].isAlive = false
								hitAlien = true
								score = score + 10
								break
							}
						}
						if hitAlien {
							break
						}
					}

					if !hitAlien && bullets[i].Y > 0 {
						activeBullets = append(activeBullets, bullets[i])
					}
				}
			}
			bullets = activeBullets

			for i := range alienBullets {
				alienBullets[i].Y++
				hitPlayer := false
				if alienBullets[i].Y < termHeight {
					if alienBullets[i].isActive && alienBullets[i].X == player.X && alienBullets[i].Y == player.Y {
						lives--
						hitPlayer = true
						if lives <= 0 {
							gameState = "lost"
						}
					}

					if !hitPlayer && alienBullets[i].Y < termHeight {
						activeAlienBullets = append(activeAlienBullets, alienBullets[i])
					}
				}

			}
			alienBullets = activeAlienBullets
			for r := range aliens {
				for c := range aliens[r] {
					if aliens[r][c].isAlive {
						screen.SetContent(aliens[r][c].X, aliens[r][c].Y, alienRune, nil, alienStyle)

					}
				}
			}
			screen.SetContent(player.X, player.Y, playerRune, nil, tcell.StyleDefault.Foreground(tcell.ColorGreen))
			for _, b := range bullets {
				screen.SetContent(b.X, b.Y, '|', nil, bulletStyle)

			}
			for _, b := range alienBullets {
				screen.SetContent(b.X, b.Y, '|', nil, bulletStyle)
			}
			for r := range aliens {
				for c := range aliens[r] {
					if aliens[r][c].isAlive {
						if aliens[r][c].Y <= 0 {
							os.Exit(1)
						}
					}
				}
			}

			livesText := fmt.Sprintf("Lives : %d", lives)
			scoreText := fmt.Sprintf("Score : %d", score)
			printString(screen, 1, 0, tcell.StyleDefault, livesText)
			printString(screen, 2+len(livesText), 0, tcell.StyleDefault.Foreground(tcell.ColorGreen), scoreText)
			won, _ := checkAlienDeaths(aliens)
			if won {
				gameState = "won"
			}

			if gameState == "lost" {
				message := "GAME OVER"
				messageX := (termWidth / 2) - len(message)
				messageY := termHeight / 2
				printString(screen, messageX, messageY, alienStyle.Bold(true), message)
			} else if gameState == "won" {
				message := "YOU WIN!"
				messageX := (termWidth / 2) - len(message)
				messageY := termHeight / 2
				printString(screen, messageX, messageY, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true), message)
			}
			screen.Show()

		case <-quitChan:
			break Loop
		}

	}
}
