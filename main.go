package main

import (
	"fmt"
	"os"
	"unsafe"
	"syscall/js"
	"github.com/xen0ne/webgl/glutil"
	"github.com/go-gl/mathgl/mgl32"
)

const vertSource = `
attribute vec4 position;
uniform mat4 Pmatrix;
uniform mat4 Mmatrix;
uniform mat4 Vmatrix;

void main() {
		gl_Position = Pmatrix * Vmatrix * position * Mmatrix;
}
`

const fragSource =`
void main() {
  gl_FragColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`

var gl js.Value


var vertz = []float32{
	-1,  1, 0,
	 1,  1, 0,
	-1, -1, 0,
	 1, -1, 0,
}

type Drawable interface {
	Draw()
}

type Rect struct {
	verts []float32
}

func NewRect(x, y, w, h int) Rect {
	
	return Rect{
		verts: []float32{
				-1,  1, 0,
				1,  1, 0,
				-1, -1, 0,
				1, -1, 0,
		},
	}
}

func main() {
	// get elements we need
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "gocanvas")
	// setup webgl
	gl = canvas.Call("getContext", "webgl")
	if gl.IsUndefined() {
		js.Global().Call("alert", "browswer does not support webgl")
		fmt.Println("we're fucked")
		os.Exit(1)
	}


	// create and BIND to vertex buffer
	vertexBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), vertexBuffer)
	gl.Call("bufferData", gl.Get("ARRAY_BUFFER"), glutil.SliceToTypedArray(vertz), gl.Get("STATIC_DRAW"))

	// SHADERS
	// init shaders
	vertShader := gl.Call("createShader", gl.Get("VERTEX_SHADER"))
	gl.Call("shaderSource", vertShader, vertSource)
	gl.Call("compileShader", vertShader)
	fragShader := gl.Call("createShader", gl.Get("FRAGMENT_SHADER"))
	gl.Call("shaderSource", fragShader, fragSource)
	gl.Call("compileShader", fragShader)
	// combine them for shader program
	shaderProgram := gl.Call("createProgram")
	gl.Call("attachShader", shaderProgram, vertShader)
	gl.Call("attachShader", shaderProgram, fragShader)
	gl.Call("linkProgram", shaderProgram)
	
	// associate shader parameters
	PositionMatrix := gl.Call("getUniformLocation", shaderProgram, "Pmatrix")
	ViewMatrix := gl.Call("getUniformLocation", shaderProgram, "Vmatrix")
	ModelMatrix := gl.Call("getUniformLocation", shaderProgram, "Mmatrix")

	// tell the shader what it's vertex position is
	gl.Call("bindBuffer", gl.Get("ARRAY_BUFFER"), vertexBuffer)
	position := gl.Call("getAttribLocation", shaderProgram, "position")
	gl.Call("vertexAttribPointer", position, 3, gl.Get("FLOAT"), false, 0, 0)
	gl.Call("enableVertexAttribArray", position)

	// use the shader
	gl.Call("useProgram", shaderProgram)

	var fov, ratio, zNear, zFar float32
	// perspective matrix
	fov = 45
	// aspect ratio, TODO: fix this lol
	ratio = 1
	// z-clipping
	zNear = 0.1
	zFar = 100.0

	// generate and apply projection matrix
	projMatrix := mgl32.Perspective(mgl32.DegToRad(fov), ratio, zNear, zFar)
	var projMatrixBuffer *[16]float32
	projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&projMatrix))
	typedProjMatrixBuffer := glutil.SliceToTypedArray([]float32((*projMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", PositionMatrix, false, typedProjMatrixBuffer)

	// Generate and apply view matrix
	viewMatrix := mgl32.Translate3D(-0.0, 0.0, -6.0)
	var viewMatrixBuffer *[16]float32
	viewMatrixBuffer = (*[16]float32)(unsafe.Pointer(&viewMatrix))
	typedViewMatrixBuffer := glutil.SliceToTypedArray([]float32((*viewMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", ViewMatrix, false, typedViewMatrixBuffer)

	// here we would apply a model matrix
	modelMatrix := mgl32.Ident4()
	var modelMatrixBuffer *[16]float32
	modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&modelMatrix))
	typedModelMatrixBuffer := glutil.SliceToTypedArray([]float32((*modelMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", ModelMatrix, false, typedModelMatrixBuffer)


	draw()

}

func draw() {
	// clear screen
	gl.Call("clearColor", 0.0, 0.0, 0.0, 1.0)
	gl.Call("clearDepth", 1.0)
	gl.Call("enable", gl.Get("DEPTH_TEST"))
	gl.Call("depthFunc", gl.Get("LEQUAL"))

	gl.Call("clear", gl.Get("COLOR_BUFFER_BIT"))
	gl.Call("clear", gl.Get("DEPTH_BUFFER_BIT"))
	// only need the vertex buffer rn
	gl.Call("drawArrays", gl.Get("TRIANGLE_STRIP"), 0, 4, glutil.SliceToTypedArray(vertz))
}
