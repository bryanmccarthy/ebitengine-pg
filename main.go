package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	resources "github.com/hajimehoshi/ebiten/v2/examples/resources/images/flappy"
)

type Game struct {
	mouseX int
	mouseY int
}

type Player struct {
	image         *ebiten.Image
	x             float64
	y             float64
	speed         float64
	anim_frame    int
	anim_cooldown int
}

type Enemy struct {
	image *ebiten.Image
	x     float64
	y     float64
	speed float64
}

type Bullet struct {
	image    *ebiten.Image
	x        float64
	y        float64
	dx       float64
	dy       float64
	speed    float64
	distance float64
}

var (
	max_distance = 100.0
)

var priestFrames = []*ebiten.Image{}

var player = Player{
	image:         nil,
	x:             0,
	y:             0,
	speed:         5.0,
	anim_frame:    0,
	anim_cooldown: 0,
}

var enemies = []Enemy{}

var bullets = []Bullet{}

func loadImage(data []byte) *ebiten.Image {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}

func (b *Bullet) calculateDirection(targetX, targetY float64) {
	dx := targetX - b.x
	dy := targetY - b.y
	distance := math.Sqrt(dx*dx + dy*dy)
	b.dx = dx / distance
	b.dy = dy / distance
}

func (g *Game) Update() error {

	ebiten.SetWindowTitle(fmt.Sprintf("Game Title | %.2ffps", ebiten.ActualFPS()))

	// Get the mouse position
	g.mouseX, g.mouseY = ebiten.CursorPosition()

	// Player movement
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		player.y -= player.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		player.y += player.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		player.x -= player.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		player.x += player.speed
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		playerHeight := player.image.Bounds().Max.Y
		playerWidth := player.image.Bounds().Max.X

		// Create a new bullet
		bullet := Bullet{
			image:    ebiten.NewImage(8, 8),
			x:        player.x + float64(playerWidth/2),
			y:        player.y + float64(playerHeight/2),
			speed:    5,
			distance: 0,
		}
		bullet.image.Fill(color.RGBA{0, 255, 0, 255})
		bullet.calculateDirection(float64(g.mouseX), float64(g.mouseY))
		bullets = append(bullets, bullet)
	}

	// Spawn enemies at random coordinates outside the screen
	if len(enemies) < 10 {
		enemy := Enemy{
			image: loadImage(resources.Gopher_png),
			x:     float64(rand.Intn(1280)),
			y:     float64(rand.Intn(720)),
			speed: 2,
		}
		enemy.image = ebiten.NewImageFromImage(enemy.image.SubImage(image.Rect(0, 0, 32, 32)).(*ebiten.Image))
		enemies = append(enemies, enemy)
	}

	// Move the enemies towards the player
	for i, enemy := range enemies {
		if enemy.x < player.x {
			enemy.x += enemy.speed
		}
		if enemy.x > player.x {
			enemy.x -= enemy.speed
		}
		if enemy.y < player.y {
			enemy.y += enemy.speed
		}
		if enemy.y > player.y {
			enemy.y -= enemy.speed
		}
		enemies[i] = enemy
	}

	// Move the bullets towards the mouse
	for i, bullet := range bullets {
		bullet.distance += bullet.speed
		bullet.x += bullet.dx * bullet.speed
		bullet.y += bullet.dy * bullet.speed
		bullets[i] = bullet
	}

	var newBullets []Bullet
	// Remove bullets that are at max distance
	for i, bullet := range bullets {
		if bullet.distance <= max_distance {
			newBullets = append(newBullets, bullets[i])
		}
	}
	bullets = newBullets

	// Check for collisions between bullets and enemies
	for i, bullet := range bullets {
		for j, enemy := range enemies {
			if bullet.x > enemy.x && bullet.x < enemy.x+32 && bullet.y > enemy.y && bullet.y < enemy.y+32 {
				enemies = append(enemies[:j], enemies[j+1:]...)
				bullets = append(bullets[:i], bullets[i+1:]...)
			}
		}
	}

	// Update the player animation
	player.anim_cooldown++
	if player.anim_cooldown >= 10 {
		player.anim_frame++
		player.anim_cooldown = 0
	}

	if player.anim_frame >= len(priestFrames) {
		player.anim_frame = 0
	}

	player.image = priestFrames[player.anim_frame]

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the player
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(player.x, player.y)
	screen.DrawImage(player.image, opts)

	// Draw the enemies
	for _, enemy := range enemies {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(enemy.x, enemy.y)
		screen.DrawImage(enemy.image, opts)
	}

	// Draw the bullets
	for _, bullet := range bullets {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(bullet.x, bullet.y)
		screen.DrawImage(bullet.image, opts)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {

	var filenames = []string{
		"resources/character_animation/priests_idle/priest1/v1/priest1_v1_1.png",
		"resources/character_animation/priests_idle/priest1/v1/priest1_v1_2.png",
		"resources/character_animation/priests_idle/priest1/v1/priest1_v1_3.png",
		"resources/character_animation/priests_idle/priest1/v1/priest1_v1_4.png",
	}

	for _, filename := range filenames {
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		img, _, err := image.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		priestFrames = append(priestFrames, ebiten.NewImageFromImage(img))

	}

	player.image = priestFrames[0]

	game := &Game{}
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Ebitengine Playground")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
