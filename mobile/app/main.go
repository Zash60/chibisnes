//go:build android
// +build android

package main

import (
	"log"
	"os"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/gl"

	"github.com/kaishuu0123/chibisnes/mobile"
)

// No global OpenGL variables needed for basic implementation

func main() {
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					onStart(a)
				case lifecycle.CrossOff:
					onStop(a)
				}
			case size.Event:
				// Handle resize if needed
			case paint.Event:
				onPaint(a)
			case touch.Event:
				onTouch(e)
			}
		}
	})
}

func onStart(a app.App) {
	// Try to load ROM from assets or create empty console
	if romData := loadRomFromAssets(); romData != nil {
		if err := mobile.Start(romData); err != "" {
			log.Printf("Failed to start with ROM: %s", err)
			mobile.Start(nil) // Start with empty console
		}
	} else {
		mobile.Start(nil) // Start with empty console
	}

	// Initialize basic OpenGL setup
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func onStop(a app.App) {
	// Cleanup if needed
}

func onPaint(a app.App) {
	// Run emulator frame
	pixels := mobile.RunFrame()

	// For now, just clear screen with a color
	// TODO: Implement proper SNES rendering
	if pixels != nil {
		// Use pixels data for rendering
		gl.ClearColor(0, 0.5, 0, 1) // Green background when emulator is running
	} else {
		gl.ClearColor(0.5, 0, 0, 1) // Red background when no ROM loaded
	}

	gl.Clear(gl.COLOR_BUFFER_BIT)
	a.Publish()
}

func onTouch(e touch.Event) {
	// Handle touch input - simplified version
	if e.Type == touch.TypeBegin || e.Type == touch.TypeMove {
		// Basic touch handling - can be expanded
		mobile.SetInput(0, true) // Example button press
	}
}

func loadRomFromAssets() []byte {
	// Try to load game.sfc from assets
	f, err := os.Open("game.sfc")
	if err != nil {
		return nil
	}
	defer f.Close()

	data := make([]byte, 1024*1024) // 1MB buffer
	n, err := f.Read(data)
	if err != nil || n == 0 {
		return nil
	}
	return data[:n]
}

// TODO: Implement proper SNES rendering with shaders when needed
