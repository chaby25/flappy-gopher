package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type bird struct {
	time int
	textures []*sdl.Texture

	y,speed float64
}

const (
	gravity = 0.1
	jumpSpeed = -5
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

func (bird *bird) paint (renderer *sdl.Renderer) error {
	bird.time++
	bird.y -= bird.speed
	bird.speed += gravity
	if bird.y < 0 {
		bird.speed = -bird.speed
		bird.y = 0
	}

	rect := &sdl.Rect{X:10,Y: (600 - int32(bird.y)) - 43/2,W:50,H:43}
	i := bird.time/10 % len(bird.textures)
	if err := renderer.Copy(bird.textures[i] , nil, rect); err != nil {
		return fmt.Errorf("Could not copy Background: %v", err)
	}
	return nil
}

func (bird *bird) destroy() {
	for _,bird := range bird.textures {
		bird.Destroy()
	}
}

func (bird *bird) jump() {
	bird.speed = jumpSpeed
}