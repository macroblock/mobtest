package main

import (
	"C"
	"fmt"

	"github.com/go-gl/mathgl/mgl32"

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
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					switch t.Keysym.Sym {
					case sdl.K_UP:
						ctx.view = ctx.view.Mul4(mgl32.Translate3D(0.0, 0.0, 0.1))
					case sdl.K_DOWN:
						ctx.view = ctx.view.Mul4(mgl32.Translate3D(0.0, 0.0, -0.1))
					case sdl.K_LEFT:
						ctx.model = ctx.model.Mul4(mgl32.HomogRotate3D(-0.1, mgl32.Vec3{0, 1, 0}))
					case sdl.K_RIGHT:
						ctx.model = ctx.model.Mul4(mgl32.HomogRotate3D(0.1, mgl32.Vec3{0, 1, 0}))
					}
				}
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
        layout(location=0) uniform mat4 aModel;
        layout(location=1) uniform mat4 aView;
        layout(location=2) uniform mat4 aProj;
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
		// 1.0, 0.0, 0.0,
		// 0.0, 1.0, 0.0,
		// 0.0, 0.0, 1.0,
	}
	_ = vertices

	elements := []uint32{2, 1, 0}
	_ = elements

	color := []float32{
		1.0, 0.0, 0.0, 1.0,
		0.0, 1.0, 0.0, 1.0,
		0.0, 0.0, 1.0, 1.0,
	}
	_ = color

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

	// gl.GenVertexArrays(1, &ctx.vao)
	// gl.BindVertexArray(ctx.vao)

	// dataBuf := NewBuffer()
	// dataBuf.Bind()
	// dataBuf.Data(vertices)
	// gl.EnableVertexAttribArray(0)
	// gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// colorBuf := NewBuffer()
	// colorBuf.Bind()
	// colorBuf.Data(color)
	// gl.EnableVertexAttribArray(1)
	// gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))

	// ctx.elements = elements
	// elemBuf := NewBuffer()
	// elemBuf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	// elemBuf.Data(elements)

	gl.GenVertexArrays(1, &ctx.vao)
	gl.BindVertexArray(ctx.vao)

	v, c, e, arr := makeCube(1.0)
	ctx.elements = arr
	v.Bind()
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	fmt.Printf("v %v\n", v)

	c.Bind()
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, gl.PtrOffset(0))
	fmt.Printf("c %v\n", c)

	e.Bind(gl.ELEMENT_ARRAY_BUFFER)
	fmt.Printf("e %v\n", e)

	ctx.view = mgl32.Translate3D(0.0, 0.0, -8)
	ctx.model = mgl32.Ident4()

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LEQUAL)
}

func onStop(ctx *TContext) {
	gl.DeleteProgram(program)
}

var (
	angle float32
)

func onDraw(ctx *TContext, w, h int32) {
	gl.Viewport(0, 0, w, h)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// angle += 0.01
	// ctx.model = mgl32.HomogRotate3D(angle, mgl32.Vec3{1, 1, 1})
	gl.UniformMatrix4fv(0, 1, false, &ctx.model[0])

	// ctx.view = mgl32.Translate3D(0.0, 0.0, -8)
	gl.UniformMatrix4fv(1, 1, false, &ctx.view[0])

	aspect := float32(w) / float32(h)
	half := float32(0.5)
	ctx.proj = mgl32.Frustum(
		-half, half,
		-half/aspect, half/aspect,
		2*half, half+10.0)
	gl.UniformMatrix4fv(2, 1, false, &ctx.proj[0])

	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 4*2*3, gl.UNSIGNED_SHORT, unsafe.Pointer(uintptr(0))) //gl.Ptr(&ctx.elements[0]))
	// gl.DrawElements(gl.TRIANGLES, 12, gl.UNSIGNED_SHORT, gl.Ptr(&ctx.elements[0]))
	// gl.BindVertexArray(ctx.vao)
	// gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
