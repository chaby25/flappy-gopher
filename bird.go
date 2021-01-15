package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type bird struct {
	mu sync.RWMutex

	time int
	textures []*sdl.Texture
	dead bool
	y,speed float64
}

const (
	gravity = 0.1
	jumpSpeed = 5
)

func newBird(renderer *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1 ; i <= 4; i++ {
		path:= fmt.Sprintf("resources/images/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(renderer,path)
		if err != nil {
			return nil,fmt.Errorf("Could not load bird image %v", err)
		}
		textures = append(textures,texture)
	}

	return &bird{textures: textures, y:300, speed: 0},nil
}

func (bird *bird) isDead() bool {
	bird.mu.RLock()
	defer bird.mu.RUnlock()

	return bird.dead
}
func (bird *bird) update ()  {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.time++
	bird.y -= bird.speed
	if bird.y < 0 {
		bird.speed = -bird.speed
		bird.y = 0
		bird.dead = true
	}
	bird.speed += gravity
}

func (bird *bird) paint (renderer *sdl.Renderer) error {
	bird.mu.RLock()
	defer bird.mu.RUnlock()

	rect := &sdl.Rect{X:10,Y: (600 - int32(bird.y)) - 43/2,W:50,H:43}
	i := bird.time/10 % len(bird.textures)
	if err := renderer.Copy(bird.textures[i] , nil, rect); err != nil {
		return fmt.Errorf("Could not copy Background: %v", err)
	}
	return nil
}

func (bird *bird) restart() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.y = 300
	bird.speed = 0
	bird.dead = false
}

func (bird *bird) destroy() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	for _,bird := range bird.textures {
		bird.Destroy()
	}
}

func (bird *bird) jump() {
	bird.mu.Lock()
	defer bird.mu.Unlock()

	bird.speed = -jumpSpeed
}