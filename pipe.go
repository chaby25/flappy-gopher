package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
	"sync"
	"time"
)


type pipes struct {
	mu sync.RWMutex
	texture *sdl.Texture
	speed int32

	pipes []*pipe
}

type pipe struct {
	mu sync.RWMutex

	x int32
	height int32
	width int32
	inverted bool
}

func newPipes(renderer *sdl.Renderer) (*pipes,error) {
	texture,err := img.LoadTexture(renderer,"resources/images/pipe.png")
	if err != nil {
		return nil,fmt.Errorf("Could not load pipe image: %v",err)
	}

	pipes := &pipes{
		texture: texture,
		speed: 4,
	}
	go func() {
		for {
			pipes.mu.Lock()
			pipes.pipes = append(pipes.pipes, newPipe())
			pipes.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()
	return pipes,nil
}

func (pipes *pipes) paint(renderer *sdl.Renderer) error {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	for _, pipe := range pipes.pipes {
		if err := pipe.paint(renderer,pipes.texture) ; err != nil {
			return err
		}
	}
	return nil
}
func (pipes *pipes) restart() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	pipes.pipes = nil
}
func (pipes *pipes) update() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	var remaining []*pipe

	for _, pipe := range pipes.pipes {
		pipe.x -= pipes.speed
		if pipe.x+pipe.width > 0 {
			remaining = append(remaining,pipe)
		}
	}
	pipes.pipes = remaining
}
func (pipes *pipes) destroy() {
	pipes.mu.Lock()
	defer pipes.mu.Unlock()

	pipes.texture.Destroy()
}

func (pipes *pipes) touch (bird *bird) {
	pipes.mu.RLock()
	defer pipes.mu.RUnlock()

	for _, pipe := range pipes.pipes {
		pipe.touch(bird)
	}
}

func (pipe *pipe) touch(bird *bird) {
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	bird.touch(pipe)
}

func newPipe() (*pipe) {
	return &pipe{
		x: 800,
		height:100 + int32(rand.Intn(300)),
		width: 50,
		inverted: rand.Float32() > 0.5,
	}
}

func (pipe *pipe) paint(renderer *sdl.Renderer, texture *sdl.Texture) error {
	pipe.mu.RLock()
	defer pipe.mu.RUnlock()

	rect := &sdl.Rect{
		X: pipe.x,
		Y: 600 - pipe.height,
		W: pipe.width,
		H: pipe.height,
	}
	flip := sdl.FLIP_NONE
	if pipe.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}

	if err := renderer.CopyEx(texture,nil,rect,0,nil,flip); err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}

	return nil

}
func (pipe *pipe) restart() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	pipe.x = 400
}

