package main

import "C"
import (
	"fmt"
	"runtime"
	"strings"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"

	gl "github.com/go-gl/gl/v3.1/gles2"
	"github.com/go-gl/mathgl/mgl32"
)

// TContext -
type TContext struct {
	window   *sdl.Window
	context  sdl.GLContext
	elements []uint32
	vao      uint32
	model    mgl32.Mat4
	view     mgl32.Mat4
	proj     mgl32.Mat4
}

func logf(format string, args ...interface{}) {
	fmt.Printf(format, args)
}

func logErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args)
}

func logPanic(format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args))
}

func esGetString(val uint32) string {
	x := gl.GetString(val)
	p := (*C.char)(unsafe.Pointer(x))
	return C.GoString(p)
}

func sdlVersion() *sdl.Version {
	v := &sdl.Version{}
	sdl.GetVersion(v)
	return v
}

// var glob = struct {
// 	window  *sdl.Window
// 	context sdl.GLContext
// }{}

func fwClose(ctx **TContext) {
	c := *ctx
	if c.context != nil {
		sdl.GLDeleteContext(c.context)
		c.context = nil
	}
	if c.window != nil {
		c.window.Destroy()
		c.window = nil
	}
	sdl.Quit()
	*ctx = nil
}

func fwInit(title string, w, h int) (*TContext, error) {
	ctx := &TContext{}
	err := error(nil)
	runtime.LockOSThread()
	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, logErrorf("sdl.Init: %v", err)
	}

	err = sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	if err != nil {
		return nil, logErrorf("sdl.GLSetAttribute: %v", err)
	}

	ctx.window, err = sdl.CreateWindow(title,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(w), int32(h),
		sdl.WINDOW_RESIZABLE|sdl.WINDOW_OPENGL)
	if err != nil {
		return nil, logErrorf("sdl.CreateWindow: %v", err)
	}

	ctx.context, err = ctx.window.GLCreateContext()
	if err != nil {
		return nil, logErrorf("sdl.GLCreateContext: %v", err)
	}

	err = gl.Init()
	if err != nil {
		return nil, logErrorf("gles.Init: %v", err)
	}
	return ctx, nil
}

func compileShader(shaderType uint32, src string) (uint32, error) {

	shader := gl.CreateShader(shaderType)

	if shader == 0 {
		// ???
		return 0, fmt.Errorf("unable to create shader: %v", gl.GetError())
	}

	csrc, free := gl.Strs(src)
	gl.ShaderSource(shader, 1, csrc, nil)
	free()
	gl.CompileShader(shader)

	status := int32(0)
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		logLen := int32(0)
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(shader, logLen, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", src, log)
	}
	return shader, nil
}

func newProgram(vSrc, fSrc string) (uint32, error) {
	vShader, err := compileShader(gl.VERTEX_SHADER, vSrc)
	if err != nil {
		return 0, err
	}

	fShader, err := compileShader(gl.FRAGMENT_SHADER, fSrc)
	if err != nil {
		return 0, err
	}

	prog := gl.CreateProgram()

	gl.AttachShader(prog, vShader)
	gl.AttachShader(prog, fShader)
	gl.LinkProgram(prog)

	status := int32(0)
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		logLen := int32(0)
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(prog, logLen, nil, gl.Str(log))

		gl.DeleteProgram(prog) // ???
		return 0, fmt.Errorf("failed to link program %v: %v", prog, log)
	}
	gl.DeleteShader(vShader)
	gl.DeleteShader(fShader)

	return prog, nil
}

func makeCube(size float32) (vert *TBuffer, color *TBuffer, elem *TBuffer, arr []uint32) {
	l := size / 2.0
	v := float32(1.0)
	z := float32(0.0)
	vertices := []float32{
		// 0.01, 0.5, 0.0,
		// -0.5, -0.5, 0.0,
		// 0.5, -0.5, 0.0,
		-l, l * 1.0, z,
		l, l, z,
		l, -l, z,
		-l, -l, z,
	}
	colors := []float32{
		// 1.0, 0.0, 0.0, 1.0,
		// 0.0, 1.0, 0.0, 1.0,
		// 0.0, 0.0, 1.0, 1.0,
		z, v, v, v,
		v, v, v, v,
		v, z, v, v,
		z, z, v, v,
		z, v, z, v,
		v, v, z, v,
		v, z, z, v,
		z, z, z, v,
	}
	elems := []uint32{
		// 2, 1, 0,
		0, 3, 1,
		3, 2, 1,
		5, 1, 2,
		5, 2, 6,
		5, 6, 7,
		5, 7, 4,
		5, 4, 0,
		5, 0, 1,
		3, 2, 1,
		3, 1, 0,
		3, 0, 4,
		3, 4, 7,
		3, 7, 6,
		3, 6, 2,
	}
	arr = elems

	vert = NewBuffer()
	vert.Bind()
	vert.Data(vertices)
	color = NewBuffer()
	color.Bind()
	color.Data(colors)
	elem = NewBuffer()
	elem.Bind(gl.ELEMENT_ARRAY_BUFFER)
	elem.Data(elems)
	return
}
