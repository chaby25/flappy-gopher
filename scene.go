package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

type scene struct {
	time int
	bg *sdl.Texture
	bird *bird
}

func newScene(renderer *sdl.Renderer) (*scene, error) {

	bg,err := img.LoadTexture(renderer,"resources/images/background.png")
	if err != nil {
		return nil,fmt.Errorf("Could not load background image")
	}

	bird, err := newBird(renderer)
	if err != nil {
		return nil,err
	}



	return &scene{bg:bg, bird:bird} ,nil
}

func (scene *scene) run(events <-chan sdl.Event, renderer *sdl.Renderer) chan error {
	errc := make (chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10*time.Millisecond)
		for {
			select {
			case e:= <-events:
				log.Printf("Event: %T", e)
			case <-tick:
					if err := scene.paint(renderer) ; err != nil {
						errc <- err
				}
			}
		}
	}()

	return errc
}
func (scene *scene) paint(renderer *sdl.Renderer) error {
	scene.time++
	renderer.Clear()

	if err := renderer.Copy(scene.bg , nil, nil); err != nil {
		return fmt.Errorf("Could not copy Background: %v", err)
	}

	if err := scene.bird.paint(renderer) ; err != nil {
		return fmt.Errorf("Error rendering bird: %v",err)
	}
	renderer.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
}