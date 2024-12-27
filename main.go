package main

import (
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	dirUp           = Point{x: 0, y: -1}
	dirDown         = Point{x: 0, y: 1}
	dirLeft         = Point{x: -1, y: 0}
	dirRight        = Point{x: 1, y: 0}
	mplusFaceSource *text.GoTextFaceSource
)

const (
	gameSpeed    = time.Second / 6
	screenWith   = 640
	screenHeight = 480
	gridSize     = 20
)

type Point struct {
	x, y int
}

type Game struct {
	snake      []Point
	direction  Point
	lastUpdate time.Time
	food       Point
	gameOver   bool
}

func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.direction = dirUp
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.direction = dirDown
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.direction = dirRight
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.direction = dirLeft
	}

	if time.Since(g.lastUpdate) < gameSpeed {
		return nil
	}
	g.lastUpdate = time.Now()
	g.updateSnake()
	return nil
}

func (g *Game) updateSnake() {
	head := (g.snake)[0]

	newHead := Point{head.x + g.direction.x, head.y + g.direction.y}

	if g.isBadCollision(newHead) {
		g.gameOver = true
	}

	if newHead == g.food {
		g.snake = append([]Point{newHead}, g.snake...)
		g.spawnFood()
	} else {
		g.snake = append(
			[]Point{newHead},
			(g.snake)[:len(g.snake)-1]...,
		)
	}
}

func (g Game) isBadCollision(p Point) bool {
	if p.x <= 0 || p.y < 0 || p.x >= screenWith/gridSize || p.y >= screenHeight/gridSize {
		return true
	}

	for _, sp := range g.snake {
		if sp == p {
			return true
		}
	}

	return false
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameOver {
		face := &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   48,
		}

		t := "Game Over!"

		w, h := text.Measure(t, face, face.Size)
		op := &text.DrawOptions{}
		op.GeoM.Translate(screenWith/2-w/2, screenHeight/2-h/2)
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, t, face, op)
	}

	for _, p := range g.snake {
		vector.DrawFilledRect(
			screen,
			float32(p.x*gridSize),
			float32(p.y*gridSize),
			gridSize,
			gridSize,
			color.White,
			true,
		)
	}
	vector.DrawFilledRect(
		screen,
		float32(g.food.x*gridSize),
		float32(g.food.y*gridSize),
		gridSize,
		gridSize,
		color.RGBA{255, 0, 0, 255},
		true,
	)
}

func (g *Game) Layout(outsideWith, outsideHeight int) (int, int) {
	return screenWith, screenHeight
}

func (g *Game) spawnFood() {
	g.food = Point{x: rand.Intn(screenWith / gridSize), y: rand.Intn(screenHeight / gridSize)}
}

func main() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	mplusFaceSource = s

	g := &Game{
		snake: []Point{
			{
				x: screenWith / gridSize / 2,
				y: screenHeight / gridSize / 2,
			},
		},
		direction: Point{x: 1, y: 0},
	}

	g.spawnFood()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Snake")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
