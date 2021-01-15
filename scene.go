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
	pipes *pipes
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

	pipes, err := newPipes(renderer)
	if err != nil {
		return nil,err
	}



	return &scene{bg:bg, bird:bird, pipes:pipes} ,nil
}

func (scene *scene) run(events <-chan sdl.Event, renderer *sdl.Renderer) chan error {
	errc := make (chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10*time.Millisecond)
		done := false
		for !done {
			select {
			case e:= <-events:
				done = scene.handleEvent(e)
			case <-tick:
				scene.update()
				if scene.bird.isDead() {
					drawTitle(renderer,"Game Over")
					time.Sleep(time.Second)
					scene.restart()
				}
				if err := scene.paint(renderer) ; err != nil {
						errc <- err
				}
			}
		}
	}()

	return errc
}

func (scene *scene) restart() {
	scene.bird.restart()
	scene.pipes.restart()

}
func (scene *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
		scene.bird.jump()
	default:
			log.Printf("Uknown Event: %T",event)
	}
	return false
}

func (scene *scene) update()  {
	scene.bird.update()
	scene.pipes.update()
	scene.pipes.touch(scene.bird)
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

	if err := scene.pipes.paint(renderer) ; err != nil {
		return fmt.Errorf("Error rendering pipe: %v",err)
	}

	renderer.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}