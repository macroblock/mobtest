package main

import (
	"C"
	"fmt"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/veandco/go-sdl2/sdl"
)
import "unsafe"

var (
	xpos, ypos float32
	program    uint32
)

func main() {
	ctx, err := fwInit("gl", 800, 600)
	defer fwClose(&ctx)
	if err != nil {
		logf("%v", err)
		return
	}

	fmt.Printf("sdl\n")
	fmt.Printf("  version: %v\n", sdlVersion())

	fmt.Printf("opengl\n")
	fmt.Printf("  version       : %v\n", esGetString(gl.VERSION))
	fmt.Printf("  shader version: %v\n", esGetString(gl.SHADING_LANGUAGE_VERSION))
	fmt.Printf("  vendor        : %v\n", esGetString(gl.VENDOR))
	fmt.Printf("  renderer      : %v\n", esGetString(gl.RENDERER))
	// fmt.Printf("  opengl extensions: \n%v\n", strings.Replace(esGetString(gl.EXTENSIONS), " ", "\n", -1))

	onStart(ctx)
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				onStop(ctx)
				running = false
			case *sdl.MouseMotionEvent:
				xpos = float32(t.X)
				ypos = float32(t.Y)
			case *sdl.WindowEvent:
				w, h := ctx.window.GetSize()
				onDraw(ctx, w, h)
				ctx.window.GLSwap()
				// sdl.PushEvent(&sdl.WindowEvent{})
			}
		}
	}
}

func onStart(ctx *TContext) {
	vShader := `#version 300 es
        layout(location=0) in vec4 aPosition;
        // layout(location=1) in vec4 aColor;
        // out vec4 vColor;
        void main() {
            gl_Position = aPosition;
            // vColor = aColor;
        }
    ` + "\x00"
	fShader := `#version 300 es
        precision mediump float;
        in vec4 vColor;
        out vec4 outColor;
        void main() {
            outColor = vColor;

        }
    ` + "\x00"
	vertices := []float32{
		0.01, 0.5, 0.0,
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
	}
	elements := []uint32{2, 1, 0}
	// color := []float32{0.0, 0.5, 0.5, 1.0}

	err := error(nil)
	program, err = newProgram(vShader, fShader)
	if err != nil {
		logPanic("%v", err)
	}

	// testdata := []byte{0, 1, 2, 3, 4}
	// buf := &TBuffer{}
	// buf.Data(vertices)
	// buf.Data(elements)
	// buf.Data(testdata)

	gl.ClearColor(0.0, 0.5, 0.0, 1.0)

	gl.UseProgram(program)

	// gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, unsafe.Pointer(&vertices[0]))
	// gl.EnableVertexAttribArray(0)
	// gl.VertexAttrib4fv(1, &color[0])

	gl.GenVertexArrays(1, &ctx.vao)
	gl.BindVertexArray(ctx.vao)

	dataBuf := NewBuffer()
	dataBuf.Bind()
	dataBuf.Data(vertices, vertices)
	// vbo(vertices)

	// var dataBuf uint32
	// gl.GenBuffers(1, &dataBuf)
	// gl.BindBuffer(gl.ARRAY_BUFFER, dataBuf)
	// gl.BufferData(
	// 	gl.ARRAY_BUFFER,
	// 	len(vertices)*4,
	// 	gl.Ptr(&vertices[0]),
	// 	gl.STATIC_DRAW)

	ctx.elements = elements
	// elemBuf := NewBuffer()
	// elemBuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	// elemBuf.Data(elements)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
}

func onStop(ctx *TContext) {
	gl.DeleteProgram(program)
}

func onDraw(ctx *TContext, w, h int32) {
	gl.Viewport(0, 0, w, h)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, unsafe.Pointer(&ctx.elements[0]))
	// gl.BindVertexArray(ctx.vao)
	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
