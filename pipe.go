package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

type pipe struct {
	mu sync.RWMutex

	texture *sdl.Texture
	x int32
	height int32
	width int32
	speed int32
	inverted bool
}

func newPipe(renderer *sdl.Renderer) (*pipe,error) {
	texture,err := img.LoadTexture(renderer,"resources/images/pipe.png")
	if err != nil {
		return nil,fmt.Errorf("Could not load pipe image: %v",err)
	}

	return &pipe{
		texture: texture,
		x: 400,
		height:300,
		width: 50,
		speed: 1,
		inverted: true,
	},nil
}

func (pipe *pipe) paint(renderer *sdl.Renderer) error {
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

	if err := renderer.CopyEx(pipe.texture,nil,rect,0,nil,flip); err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}

	return nil

}
func (pipe *pipe) restart() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	pipe.x = 400
}
func (pipe *pipe) update() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()
	pipe.x -= pipe.speed

}
func (pipe *pipe) destroy() {
	pipe.mu.Lock()
	defer pipe.mu.Unlock()

	pipe.texture.Destroy()
}