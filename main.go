package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func mandelbrot(cx, cy float64, maxIter int) float64 {
	x, y := 0.0, 0.0
	iteration := 0

	for x*x+y*y <= 4 && iteration < maxIter {
		xTemp := x*x - y*y + cx
		y = 2*x*y + cy
		x = xTemp
		iteration++
	}

	if iteration < maxIter {
		logZn := math.Log(x*x+y*y) / 2
		return float64(iteration) + 1 - math.Log(logZn)/math.Log(2)
	}
	return float64(maxIter)
}

// colour mapping from: https://stackoverflow.com/questions/16500656/which-color-gradient-is-used-to-color-mandelbrot-in-wikipedia
var colorMapping = []color.RGBA{
	{66, 30, 15, 255},
	{25, 7, 26, 255},
	{9, 1, 47, 255},
	{4, 4, 73, 255},
	{0, 7, 100, 255},
	{12, 44, 138, 255},
	{24, 82, 177, 255},
	{57, 125, 209, 255},
	{134, 181, 229, 255},
	{211, 236, 248, 255},
	{241, 233, 191, 255},
	{248, 201, 95, 255},
	{255, 170, 0, 255},
	{204, 128, 0, 255},
	{153, 87, 0, 255},
	{106, 52, 3, 255},
}

func getColor(iterations, maxIter int) color.RGBA {
	if iterations < maxIter && iterations > 0 {
		i := iterations % len(colorMapping)
		return colorMapping[i]
	}
	return color.RGBA{}
}

type Game struct {
	minX, maxX, minY, maxY float64
	centerX, centerY       float64
	zoom                   float64
	zoomSpeed              float64
	lastUpdate             time.Time
}

func (g *Game) Update() error {
	now := time.Now()
	elapsed := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// sidebar interaction
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x < 100 {
			if y >= 70 && y <= 270 {
				g.zoomSpeed = (float64(y-70) / 200) * 0.5 // Max speed of 0.5
			}
		}
	}

	g.zoom *= math.Pow(1+g.zoomSpeed, elapsed)
	g.zoom = math.Max(1, math.Min(g.zoom, 1e15))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	maxIter := 100

	width := (g.maxX - g.minX) / g.zoom
	height := (g.maxY - g.minY) / g.zoom
	minX := g.centerX - width/2
	maxX := g.centerX + width/2
	minY := g.centerY - height/2
	maxY := g.centerY + height/2

	// calc mandelbrot set for each pixel
	for y := 0; y < screen.Bounds().Dy(); y++ {
		for x := 0; x < screen.Bounds().Dx(); x++ {
			cx := minX + (maxX-minX)*float64(x)/float64(screen.Bounds().Dx())
			cy := minY + (maxY-minY)*float64(y)/float64(screen.Bounds().Dy())

			iterations := mandelbrot(cx, cy, maxIter)
			clr := getColor(int(iterations), maxIter)

			vector.DrawFilledRect(screen, float32(x), float32(y), 1, 1, clr, false)
		}
	}

	drawSidebar(screen, g)
	drawInfo(screen, g.zoomSpeed, g.zoom, g.centerX, g.centerY)
}

func drawSidebar(screen *ebiten.Image, g *Game) {
	sidebarWidth := 100
	sidebarColor := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	sidebarRect := ebiten.NewImage(sidebarWidth, screen.Bounds().Dy())
	sidebarRect.Fill(sidebarColor)

	screen.DrawImage(sidebarRect, nil)

	zoomSpeedX := 10
	zoomSpeedY := 70
	zoomSpeedHeight := 200

	vector.DrawFilledRect(screen, float32(zoomSpeedX), float32(zoomSpeedY), 10, float32(zoomSpeedHeight), color.RGBA{200, 200, 200, 255}, false)

	// Draw the current zoom speed as a filled rectangle
	currentZoomSpeedY := zoomSpeedY + int((g.zoomSpeed/0.5)*float64(zoomSpeedHeight))
	vector.DrawFilledRect(screen, float32(zoomSpeedX), float32(currentZoomSpeedY-5), 10, 10, color.RGBA{255, 0, 0, 255}, false)
}

func drawInfo(screen *ebiten.Image, zoomSpeed, zoomLevel, centerX, centerY float64) {
	myFont := basicfont.Face7x13

	speedContent := fmt.Sprintf("Zoom Speed: %.3f", zoomSpeed)
	text.Draw(screen, speedContent, myFont, 10, 20, color.White)

	levelContent := fmt.Sprintf("Zoom Level: %.2f", zoomLevel)
	text.Draw(screen, levelContent, myFont, 10, 40, color.White)

	centerContent := fmt.Sprintf("Center: (%.6f, %.6f)", centerX, centerY)
	text.Draw(screen, centerContent, myFont, 10, 60, color.White)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	game := &Game{
		minX: -2.5,
		maxX: 1.0,
		minY: -1.5,
		maxY: 1.5,
		/* Center on Seahorse Valley
		http://www.mrob.com/pub/muency/seahorsevalley.html
		*/
		centerX:    -0.7445398603559083806,
		centerY:    0.1217237738944248242,
		zoom:       0.0,  // Initial zoom level
		zoomSpeed:  0.01, // Initial zoom speed
		lastUpdate: time.Now(),
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Mandelbrot Set")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
