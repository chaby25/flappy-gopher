package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"runtime"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

func run() error{
	sdl.SetHint(sdl.HINT_RENDER_VSYNC, "1")

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("Could not initialize SDL: %v",err)

	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("Could not initialize TTF: %v",err)
	}
	defer ttf.Quit()

	window, renderer, err :=  sdl.CreateWindowAndRenderer(800,600,sdl.WINDOW_RESIZABLE)
	if err != nil {
		return fmt.Errorf("Could not create Window: %v", err)
	}
	defer window.Destroy()

	if err := drawTitle(renderer); err != nil {
		return fmt.Errorf("Could not draw title: %v",err)
	}
	time.Sleep(1*time.Second)

	scene , err := newScene(renderer)

	if  err != nil {
		return fmt.Errorf("Could not draw scene: %v",err)
	}
	defer scene.destroy()

	events := make(chan sdl.Event)
	errc := scene.run(events,renderer)

	runtime.LockOSThread()
	for {
		select {
			case events <- sdl.WaitEvent():
			case err := <-errc:
				return err
		}
	}

//	return <-scene.run(events,renderer)
}

func drawTitle(renderer *sdl.Renderer) error {
	renderer.Clear()
	font, err := ttf.OpenFont("resources/fonts/Flappy.ttf",20)
	if err != nil {
		return fmt.Errorf("Could not load font: %v", err)
	}
	defer font.Close()
	surface, err := font.RenderUTF8Solid("Flappy Gopher",sdl.Color {
		R: 255,
		G: 100,
		B: 0,
		A: 255})
	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)

	if err != nil {
		return fmt.Errorf("Could not create texture: %v" , err)
	}
	defer texture.Destroy()

	if err := renderer.Copy(texture,nil,nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v" , err)
	}

	renderer.Present()


	return nil
}
