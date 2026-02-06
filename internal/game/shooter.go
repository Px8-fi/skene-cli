package game

import (
	"math/rand"
	"skene-terminal-v2/internal/tui/styles"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Entity types
type EntityType int

const (
	EntityPlayer EntityType = iota
	EntityEnemy
	EntityBullet
	EntityPowerUp
)

// Entity represents a game entity
type Entity struct {
	Type   EntityType
	X, Y   int
	Width  int
	Height int
	Alive  bool
	Sprite string
}

// Game represents the space shooter game
type Game struct {
	width        int
	height       int
	player       *Entity
	enemies      []*Entity
	bullets      []*Entity
	powerUps     []*Entity
	score        int
	lives        int
	level        int
	gameOver     bool
	paused       bool
	lastSpawn    time.Time
	spawnRate    time.Duration
}

// NewGame creates a new game instance
func NewGame(width, height int) *Game {
	g := &Game{
		width:     width,
		height:    height,
		enemies:   make([]*Entity, 0),
		bullets:   make([]*Entity, 0),
		powerUps:  make([]*Entity, 0),
		score:     0,
		lives:     3,
		level:     1,
		gameOver:  false,
		paused:    false,
		lastSpawn: time.Now(),
		spawnRate: 800 * time.Millisecond,
	}

	// Create player
	g.player = &Entity{
		Type:   EntityPlayer,
		X:      width / 2,
		Y:      height - 3,
		Width:  3,
		Height: 2,
		Alive:  true,
		Sprite: " ▲ \n/█\\",
	}

	return g
}

// SetSize updates game dimensions
func (g *Game) SetSize(width, height int) {
	g.width = width
	g.height = height
	
	// Reposition player
	if g.player != nil {
		g.player.Y = height - 3
		if g.player.X > width-3 {
			g.player.X = width - 3
		}
	}
}

// Update game state
func (g *Game) Update() {
	if g.gameOver || g.paused {
		return
	}

	// Move bullets
	for _, b := range g.bullets {
		if b.Alive {
			b.Y--
			if b.Y < 0 {
				b.Alive = false
			}
		}
	}

	// Move enemies
	for _, e := range g.enemies {
		if e.Alive {
			e.Y++
			if e.Y > g.height {
				e.Alive = false
			}
		}
	}

	// Check collisions
	g.checkCollisions()

	// Spawn new enemies
	if time.Since(g.lastSpawn) > g.spawnRate {
		g.spawnEnemy()
		g.lastSpawn = time.Now()
	}

	// Clean up dead entities
	g.cleanup()

	// Check game over
	if g.lives <= 0 {
		g.gameOver = true
	}
}

// MoveLeft moves player left
func (g *Game) MoveLeft() {
	if g.player.X > 1 {
		g.player.X -= 2
	}
}

// MoveRight moves player right
func (g *Game) MoveRight() {
	if g.player.X < g.width-4 {
		g.player.X += 2
	}
}

// Shoot fires a bullet
func (g *Game) Shoot() {
	bullet := &Entity{
		Type:   EntityBullet,
		X:      g.player.X + 1,
		Y:      g.player.Y - 1,
		Width:  1,
		Height: 1,
		Alive:  true,
		Sprite: "│",
	}
	g.bullets = append(g.bullets, bullet)
}

// spawnEnemy creates a new enemy
func (g *Game) spawnEnemy() {
	enemyTypes := []struct {
		sprite string
		width  int
	}{
		{"◆", 1},
		{"▼", 1},
		{"●", 1},
		{"◈", 1},
	}

	et := enemyTypes[rand.Intn(len(enemyTypes))]

	enemy := &Entity{
		Type:   EntityEnemy,
		X:      rand.Intn(g.width - 4) + 2,
		Y:      0,
		Width:  et.width,
		Height: 1,
		Alive:  true,
		Sprite: et.sprite,
	}
	g.enemies = append(g.enemies, enemy)
}

// checkCollisions checks for collisions
func (g *Game) checkCollisions() {
	// Bullets vs Enemies
	for _, b := range g.bullets {
		if !b.Alive {
			continue
		}
		for _, e := range g.enemies {
			if !e.Alive {
				continue
			}
			if g.collides(b, e) {
				b.Alive = false
				e.Alive = false
				g.score += 100
				
				// Level up every 1000 points
				if g.score > 0 && g.score%1000 == 0 {
					g.level++
					if g.spawnRate > 300*time.Millisecond {
						g.spawnRate -= 100 * time.Millisecond
					}
				}
			}
		}
	}

	// Enemies vs Player
	for _, e := range g.enemies {
		if !e.Alive {
			continue
		}
		if g.collides(e, g.player) {
			e.Alive = false
			g.lives--
		}
	}
}

func (g *Game) collides(a, b *Entity) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

func (g *Game) cleanup() {
	// Clean bullets
	var aliveBullets []*Entity
	for _, b := range g.bullets {
		if b.Alive {
			aliveBullets = append(aliveBullets, b)
		}
	}
	g.bullets = aliveBullets

	// Clean enemies
	var aliveEnemies []*Entity
	for _, e := range g.enemies {
		if e.Alive {
			aliveEnemies = append(aliveEnemies, e)
		}
	}
	g.enemies = aliveEnemies
}

// Restart resets the game
func (g *Game) Restart() {
	g.enemies = make([]*Entity, 0)
	g.bullets = make([]*Entity, 0)
	g.powerUps = make([]*Entity, 0)
	g.score = 0
	g.lives = 3
	g.level = 1
	g.gameOver = false
	g.paused = false
	g.spawnRate = 800 * time.Millisecond
	g.player.X = g.width / 2
	g.player.Y = g.height - 3
	g.player.Alive = true
}

// TogglePause toggles pause state
func (g *Game) TogglePause() {
	g.paused = !g.paused
}

// IsGameOver returns if game is over
func (g *Game) IsGameOver() bool {
	return g.gameOver
}

// IsPaused returns if game is paused
func (g *Game) IsPaused() bool {
	return g.paused
}

// GetScore returns current score
func (g *Game) GetScore() int {
	return g.score
}

// Render draws the game
func (g *Game) Render() string {
	// Create game field
	field := make([][]rune, g.height)
	for i := range field {
		field[i] = make([]rune, g.width)
		for j := range field[i] {
			field[i][j] = ' '
		}
	}

	// Draw stars (background)
	starPositions := []struct{ x, y int }{
		{5, 3}, {15, 7}, {25, 2}, {35, 8}, {45, 4},
		{10, 12}, {20, 15}, {30, 10}, {40, 13}, {50, 6},
	}
	for _, pos := range starPositions {
		if pos.x < g.width && pos.y < g.height {
			field[pos.y][pos.x] = '·'
		}
	}

	// Draw enemies
	for _, e := range g.enemies {
		if e.Alive && e.Y >= 0 && e.Y < g.height && e.X >= 0 && e.X < g.width {
			for i, r := range e.Sprite {
				if e.X+i < g.width {
					field[e.Y][e.X+i] = r
				}
			}
		}
	}

	// Draw bullets
	for _, b := range g.bullets {
		if b.Alive && b.Y >= 0 && b.Y < g.height && b.X >= 0 && b.X < g.width {
			field[b.Y][b.X] = '│'
		}
	}

	// Draw player
	if g.player.Alive && g.player.Y >= 0 && g.player.Y < g.height {
		// Simple player drawing
		if g.player.Y-1 >= 0 && g.player.X+1 < g.width {
			field[g.player.Y-1][g.player.X+1] = '▲'
		}
		if g.player.X < g.width {
			field[g.player.Y][g.player.X] = '/'
		}
		if g.player.X+1 < g.width {
			field[g.player.Y][g.player.X+1] = '█'
		}
		if g.player.X+2 < g.width {
			field[g.player.Y][g.player.X+2] = '\\'
		}
	}

	// Convert field to string
	var lines []string
	for _, row := range field {
		lines = append(lines, string(row))
	}

	gameArea := strings.Join(lines, "\n")

	// Style the game area
	enemyStyle := lipgloss.NewStyle().Foreground(styles.GameMagenta)
	bulletStyle := lipgloss.NewStyle().Foreground(styles.GameYellow)
	playerStyle := lipgloss.NewStyle().Foreground(styles.GameCyan)
	starStyle := lipgloss.NewStyle().Foreground(styles.MidGray)

	// Apply colors (simplified - in real implementation would be per-character)
	gameArea = starStyle.Render(gameArea)

	// Header with score
	header := lipgloss.JoinHorizontal(
		lipgloss.Center,
		styles.Accent.Render("SPACE SHOOTER"),
		"  ",
		styles.Body.Render("Score: "),
		styles.Accent.Render(strings.Repeat("█", g.score/100)),
		styles.Body.Render(string(rune('0'+g.score/100))),
		"  ",
		styles.Body.Render("Lives: "),
		styles.Accent.Render(strings.Repeat("♥", g.lives)),
		"  ",
		styles.Body.Render("Level: "),
		styles.Accent.Render(string(rune('0'+g.level))),
	)

	// Game box
	gameBox := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(styles.GameCyan).
		Render(gameArea)

	// Footer
	footer := styles.Muted.Render("← → move • space shoot • p pause • esc exit")

	// Overlay for pause/game over
	var overlay string
	if g.paused {
		overlay = lipgloss.Place(
			g.width,
			g.height,
			lipgloss.Center,
			lipgloss.Center,
			styles.Box.Render(styles.Accent.Render("PAUSED\n\nPress P to continue")),
		)
	} else if g.gameOver {
		overlay = lipgloss.Place(
			g.width,
			g.height,
			lipgloss.Center,
			lipgloss.Center,
			styles.Box.Render(
				styles.Error.Render("GAME OVER\n\n")+
					styles.Body.Render("Final Score: ")+
					styles.Accent.Render(strings.Repeat("█", g.score/100)+"\n\n")+
					styles.Muted.Render("Press R to restart • ESC to exit"),
			),
		)
	}

	// Combine
	result := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		gameBox,
		"",
		footer,
	)

	if overlay != "" {
		// Overlay the message
		result = overlay
	}

	// Apply entity-specific styling (for visibility)
	_ = enemyStyle
	_ = bulletStyle
	_ = playerStyle

	return result
}

// GameTickMsg is sent for game updates
type GameTickMsg time.Time

// GameTickCmd returns a command for game ticks
func GameTickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
		return GameTickMsg(t)
	})
}
