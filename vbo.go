package main

import (
	"fmt"
	"reflect"
	"unsafe"

	gl "github.com/go-gl/gl/v3.1/gles2"
)

// TBuffer -
type TBuffer struct {
	buffer     uint32
	targetHint uint32
	usageHint  uint32
}

// NewBuffer -
func NewBuffer() *TBuffer {
	ret := &TBuffer{targetHint: gl.ARRAY_BUFFER, usageHint: gl.STATIC_DRAW}
	gl.GenBuffers(1, &ret.buffer)
	return ret
}

// Bind -
func (o *TBuffer) Bind(targetHint ...uint32) {
	if len(targetHint) > 0 {
		o.targetHint = targetHint[0]
	}
	gl.BindBuffer(o.targetHint, o.buffer)
}

// Data -
func (o *TBuffer) Data(d0 []float32, data interface{}, usageHint ...uint32) {
	if len(usageHint) > 0 {
		o.usageHint = usageHint[0]
	}

	val := reflect.ValueOf(data)
	if val.Len() == 0 {
		// ???
		return
	}
	typ := reflect.TypeOf(data)
	typeSize := typ.Elem().Size()
	sliceLen := val.Len() // for slice, arrays or chan only
	ptr := unsafe.Pointer(val.Pointer())
	// fmt.Println(sliceLen, " ", typeSize, " ", ptr)
	gl.BufferData(
		o.targetHint,
		int(sliceLen)*int(typeSize),
		ptr,
		o.usageHint)
	// gl.GenBuffers(1, &o.buffer)
	// gl.BindBuffer(o.targetHint, o.buffer)
	// gl.BufferData(
	// 	gl.ARRAY_BUFFER,
	// 	len(d0)*4,
	// 	gl.Ptr(&d0[0]),
	// 	o.usageHint)

	fmt.Printf("-> %v\n", *(*uint32)(ptr))
}

func vbo(data []float32) {
	var dataBuf uint32
	gl.GenBuffers(1, &dataBuf)
	gl.BindBuffer(gl.ARRAY_BUFFER, dataBuf)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(data)*4,
		gl.Ptr(&data[0]),
		gl.STATIC_DRAW)
}
