// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin linux windows

// An app that draws a green triangle on a red background.
//
// Note: This demo is an early preview of Go 1.5. In order to build this
// program as an Android APK using the gomobile tool.
//
// See http://godoc.org/golang.org/x/mobile/cmd/gomobile to install gomobile.
//
// Get the basic example and use gomobile to build or install it on your device.
//
//   $ go get -d golang.org/x/mobile/example/basic
//   $ gomobile build golang.org/x/mobile/example/basic # will build an APK
//
//   # plug your Android device to your computer or start an Android emulator.
//   # if you have adb installed on your machine, use gomobile install to
//   # build and deploy the APK to an Android target.
//   $ gomobile install golang.org/x/mobile/example/basic
//
// Switch to your device or emulator to start the Basic application from
// the launcher.
// You can also run the application on your desktop by running the command
// below. (Note: It currently doesn't work on Windows.)
//   $ go install golang.org/x/mobile/example/basic && basic
package main

import (
	"encoding/binary"
	"log"

	"math"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

var (
	images   *glutil.Images
	fps      *debug.FPS
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer
	buf2     gl.Buffer
	buf3     gl.Buffer

	green  float32
	touchX float32
	touchY float32
)

func main() {
	app.Main(func(a app.App) {
		var glctx gl.Context
		var sz size.Event
		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					glctx, _ = e.DrawContext.(gl.Context)
					onStart(glctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop(glctx)
					glctx = nil
				}
			case size.Event:
				sz = e
				touchX = float32(sz.WidthPx / 2)
				touchY = float32(sz.HeightPx / 2)
			case paint.Event:
				if glctx == nil || e.External {
					// As we are actively painting as fast as
					// we can (usually 60 FPS), skip any paint
					// events sent by the system.
					continue
				}

				onPaint(glctx, sz)
				a.Publish()
				// Drive the animation by preparing to paint the next frame
				// after this one is shown.
				a.Send(paint.Event{})
			case touch.Event:
				touchX = e.X
				touchY = e.Y
			}
		}
	})
}

func onStart(glctx gl.Context) {
	var err error
	program, err = glutil.CreateProgram(glctx, vertexShader, fragmentShader)
	if err != nil {
		log.Printf("error creating GL program: %v", err)
		return
	}

	buf = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
	glctx.BufferData(gl.ARRAY_BUFFER, line_vertical, gl.STATIC_DRAW)

	buf2 = glctx.CreateBuffer()
	glctx.BindBuffer(gl.ARRAY_BUFFER, buf2)
	glctx.BufferData(gl.ARRAY_BUFFER, line_horizontal, gl.STATIC_DRAW)

	buf3 = glctx.CreateBuffer()

	position = glctx.GetAttribLocation(program, "position")
	color = glctx.GetUniformLocation(program, "color")
	offset = glctx.GetUniformLocation(program, "offset")

	images = glutil.NewImages(glctx)
	fps = debug.NewFPS(images)
}

func onStop(glctx gl.Context) {
	glctx.DeleteProgram(program)
	glctx.DeleteBuffer(buf)
	fps.Release()
	images.Release()
}

func onPaint(glctx gl.Context, sz size.Event) {

	log.Println("click coords: ", touchX, touchY)

	glctx.ClearColor(0, 0, 0, 0)
	glctx.Clear(gl.COLOR_BUFFER_BIT)

	glctx.UseProgram(program)

	var current_x float32 = -0.7
	var size_x float32 = 0.25

	for i := 0; i < 7; i++ {
		glctx.Uniform4f(color, 0, 1, 0, 1)

		glctx.Uniform2f(offset, current_x, 0)
		current_x += size_x
		// glctx.Uniform2f(offset, 0.2+(0.2*float32(i)), 0)
		// glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))
		// log.Println(touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))

		glctx.BindBuffer(gl.ARRAY_BUFFER, buf)
		glctx.EnableVertexAttribArray(position)
		glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
		glctx.DrawArrays(gl.LINES, 0, vertexCount)
		glctx.DisableVertexAttribArray(position)
	}

	var current_y float32 = -0.7
	var size_y float32 = 0.25

	for i := 0; i < 7; i++ {
		glctx.Uniform4f(color, 0, 1, 0, 1)

		glctx.Uniform2f(offset, 0, current_y)
		current_y += size_y
		// glctx.Uniform2f(offset, 0.2+(0.2*float32(i)), 0)
		// glctx.Uniform2f(offset, touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))
		// log.Println(touchX/float32(sz.WidthPx), touchY/float32(sz.HeightPx))

		glctx.BindBuffer(gl.ARRAY_BUFFER, buf2)
		glctx.EnableVertexAttribArray(position)
		glctx.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
		glctx.DrawArrays(gl.LINES, 0, vertexCount)
		glctx.DisableVertexAttribArray(position)
	}

	{
		// Draw circle
		vertex_count := 30

		var radius float32 = 0.08
		var center_x float32 = 0.0
		var center_y float32 = 0.0

		buffer := make([]float32, vertex_count*2)
		idx := 0

		buffer[idx] = center_x
		idx += 1
		buffer[idx] = center_y
		idx += 1

		outerVertexCount := vertex_count - 1

		for i := 1; i < outerVertexCount; i++ {
			var percent float64 = (float64(i) / float64(outerVertexCount-1))
			var rad float64 = percent * (2 * math.Pi)

			var outer_x float32 = center_x + (radius * float32(math.Cos(rad)))
			var outer_y float32 = center_y + (radius * float32(math.Sin(rad)))

			buffer[idx] = outer_x
			idx += 1
			buffer[idx] = outer_y
			idx += 1

		}

		// var radius float64 = 30
		// var sides int = 60

		// buffer := []float32{}

		// for i := 0.0; i < 2*math.Pi; i += (2 * math.Pi / float64(sides)) {
		// 	buffer = append(buffer, float32(math.Sin(i)*radius))
		// 	buffer = append(buffer, float32(math.Cos(i)*radius))
		// }

		glctx.Uniform4f(color, 0, 0, 1, 1)
		// glctx.Uniform2f(offset, 0, 0)

		x_offset := 2.0*(touchX/float32(sz.WidthPx)) - 1.0
		y_offset := 1.0 - 2.0*(touchY/float32(sz.HeightPx))
		glctx.Uniform2f(offset, x_offset, y_offset)

		glctx.BindBuffer(gl.ARRAY_BUFFER, buf3)
		glctx.BufferData(gl.ARRAY_BUFFER, f32.Bytes(binary.LittleEndian, buffer...), gl.STATIC_DRAW)

		glctx.EnableVertexAttribArray(position)
		glctx.VertexAttribPointer(position, 2, gl.FLOAT, false, 0, 0)
		glctx.DrawArrays(gl.LINE_LOOP, 2, outerVertexCount)
		glctx.DisableVertexAttribArray(position)
	}

	fps.Draw(sz)
}

var line_vertical = f32.Bytes(binary.LittleEndian,
	0.0, 1.0, 0.0, // top left
	0.0, -1.0, 0.0, // bottom left
)

var line_horizontal = f32.Bytes(binary.LittleEndian,
	-1.0, 0.0, 0.0, // top left
	1.0, 0.0, 0.0, // bottom left
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vertexShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
	// offset comes in with x/y values between 0 and 1.
	vec4 offset4 = vec4(offset.x, offset.y, 0, 0);
	gl_Position = position + offset4;
}`

const fragmentShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
	gl_FragColor = color;
}`
