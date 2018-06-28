package main

import (
	"C"
	"fmt"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/veandco/go-sdl2/sdl"
)

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
				// w, h := ctx.window.GetSize()
				// onDraw(ctx, w, h)
				// sdl.PushEvent(&sdl.WindowEvent{})
			}
		}
		w, h := ctx.window.GetSize()
		onDraw(ctx, w, h)
		ctx.window.GLSwap()
	}
}

func onStart(ctx *TContext) {
	vShader := `#version 300 es
        #extension GL_ARB_explicit_uniform_location : enable
        layout(location=0) in vec3 aPosition;
        layout(location=1) in vec3 aColor;
        out vec3 vColor;
        layout(location=2) uniform mat4 aModel;
        layout(location=3) uniform mat4 aView;
        layout(location=4) uniform mat4 aProj;
        void main() {
            gl_Position = aProj * aView * aModel * vec4(aPosition,1);
            vColor = aColor;
        }
    ` + "\x00"
	fShader := `#version 300 es
        precision mediump float;
        in vec3 vColor;
        out vec3 outColor;
        void main() {
            outColor = vColor;
        }
    ` + "\x00"
	vertices := []float32{
		0.01, 0.5, 0.0,
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		// color
		1.0, 0.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 0.0, 1.0,
	}
	elements := []uint32{2, 1, 0}
	// color := []float32{
	// 	1.0, 0.0, 0.0, 1.0,
	// 	0.0, 1.0, 0.0, 1.0,
	// 	0.0, 0.0, 1.0, 1.0,
	// }

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
	dataBuf.Data(vertices)
	// vbo(vertices)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, gl.PtrOffset(9*4))

	// var dataBuf uint32
	// gl.GenBuffers(1, &dataBuf)
	// gl.BindBuffer(gl.ARRAY_BUFFER, dataBuf)
	// gl.BufferData(
	// 	gl.ARRAY_BUFFER,
	// 	len(vertices)*4,
	// 	gl.Ptr(&vertices[0]),
	// 	gl.STATIC_DRAW)

	ctx.elements = elements
	elemBuf := NewBuffer()
	elemBuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	elemBuf.Data(elements)

}

func onStop(ctx *TContext) {
	gl.DeleteProgram(program)
}

var (
	angle float32
	model = mgl32.Ident4()
)

func onDraw(ctx *TContext, w, h int32) {
	gl.Viewport(0, 0, w, h)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	angle += 0.01
	model = mgl32.HomogRotate3D(angle, mgl32.Vec3{1, 1, 1})
	gl.UniformMatrix4fv(2, 1, false, &model[0])

	view := mgl32.Translate3D(0.0, 0.0, -2)
	gl.UniformMatrix4fv(3, 1, false, &view[0])

	proj := mgl32.Frustum(
		-0.5, 0.5,
		-0.5, 0.5,
		0.5, 10.0)
	gl.UniformMatrix4fv(4, 1, false, &proj[0])

	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_SHORT, unsafe.Pointer(uintptr(0))) //gl.Ptr(&ctx.elements[0]))
	// gl.BindVertexArray(ctx.vao)
	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
