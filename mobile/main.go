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

var (
	console *mobile.Console
	program gl.Program
	buf     gl.Buffer
)

func main() {
	app.Main(func(a app.App) {
		for e := range a.Events() {
			switch e := app.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					onStart()
				case lifecycle.CrossOff:
					onStop()
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

func onStart() {
	// Try to load ROM from assets or create empty console
	if romData := loadRomFromAssets(); romData != nil {
		if err := mobile.Start(romData); err != "" {
			log.Printf("Failed to start with ROM: %s", err)
			mobile.Start(nil) // Start with empty console
		}
	} else {
		mobile.Start(nil) // Start with empty console
	}

	// Initialize OpenGL
	program = gl.CreateProgram()
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, vertexShaderSource)
	gl.CompileShader(vertexShader)
	gl.AttachShader(program, vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, fragmentShaderSource)
	gl.CompileShader(fragmentShader)
	gl.AttachShader(program, fragmentShader)

	gl.LinkProgram(program)
	gl.UseProgram(program)

	buf = gl.CreateBuffer()
	vertices := []float32{
		-1, -1, 0, 1,
		1, -1, 1, 1,
		-1, 1, 0, 0,
		1, 1, 1, 0,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
}

func onStop() {
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

func onPaint(a app.App) {
	if pixels := mobile.RunFrame(); pixels != nil {
		// TODO: Render pixels using OpenGL
	}

	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 16, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 16, 8)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)
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

const vertexShaderSource = `
attribute vec2 position;
attribute vec2 texCoord;
varying vec2 vTexCoord;
void main() {
    gl_Position = vec4(position, 0.0, 1.0);
    vTexCoord = texCoord;
}
` + "\x00"

const fragmentShaderSource = `
precision mediump float;
varying vec2 vTexCoord;
void main() {
    gl_FragColor = vec4(vTexCoord, 0.0, 1.0);
}
` + "\x00"
