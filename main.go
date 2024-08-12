// graphicalboredomkiller
// I coded this while I was very sick and bored
// A simple and slightly fun example of a 2d graphical program written in go (using sdl2)

package main

import (
	"math"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

// Shape types
const (
	ShapeRectFilled      = 0
	ShapeRectOutline     = 1
	ShapeTriangleFilled  = 2
	ShapeTriangleOutline = 3
)

type Shape struct {
	ShapeType      uint8
	X              int32
	Y              int32
	Size           int32
	Color          sdl.Color
	Angle          float32
	MovementVector Vector
}

type Vector struct {
	X float32
	Y float32
}

var Shapes []Shape
var TempShapes []Shape

func genTrianglePoints(x float32, y float32, length int32, angle float32) []Vector {
	Vector1 := Vector{X: float32(length), Y: 0}
	Vector2 := Vector{X: float32(length) / 2, Y: float32(math.Sqrt(math.Pow(float64(length), 2) - (math.Pow(float64(length), 2))/4))}

	if angle != 0 {
		Vector1 = Vector{
			X: Vector1.X*float32(math.Cos(float64(angle))) - Vector1.Y*float32(math.Sin(float64(angle))),
			Y: Vector1.X*float32(math.Sin(float64(angle))) + Vector1.Y*float32(math.Cos(float64(angle)))}
		Vector2 = Vector{
			X: Vector2.X*float32(math.Cos(float64(angle))) - Vector2.Y*float32(math.Sin(float64(angle))),
			Y: Vector2.X*float32(math.Sin(float64(angle))) + Vector2.Y*float32(math.Cos(float64(angle)))}
	}

	return []Vector{Vector{X: x, Y: y}, Vector{X: Vector1.X + x, Y: Vector1.Y + y}, Vector{X: Vector2.X + x, Y: Vector2.Y + y}}
}

func genRectPoints(x float32, y float32, length int32, angle float32) []Vector {
	Vector1 := Vector{X: float32(length), Y: 0}
	Vector2 := Vector{X: float32(length), Y: float32(length)}
	Vector3 := Vector{X: 0, Y: float32(length)}

	if angle != 0 {
		Vector1 = Vector{
			X: Vector1.X*float32(math.Cos(float64(angle))) - Vector1.Y*float32(math.Sin(float64(angle))),
			Y: Vector1.X*float32(math.Sin(float64(angle))) + Vector1.Y*float32(math.Cos(float64(angle)))}
		Vector2 = Vector{
			X: Vector2.X*float32(math.Cos(float64(angle))) - Vector2.Y*float32(math.Sin(float64(angle))),
			Y: Vector2.X*float32(math.Sin(float64(angle))) + Vector2.Y*float32(math.Cos(float64(angle)))}
		Vector3 = Vector{
			X: Vector3.X*float32(math.Cos(float64(angle))) - Vector3.Y*float32(math.Sin(float64(angle))),
			Y: Vector3.X*float32(math.Sin(float64(angle))) + Vector3.Y*float32(math.Cos(float64(angle)))}
	}

	return []Vector{Vector{X: x, Y: y}, Vector{X: Vector1.X + x, Y: Vector1.Y + y}, Vector{X: Vector2.X + x, Y: Vector2.Y + y}, Vector{X: Vector3.X + x, Y: Vector3.Y + y}}
}

func pointToVertex(v Vector, c sdl.Color) sdl.Vertex {
	return sdl.Vertex{Position: sdl.FPoint{X: v.X, Y: v.Y}, Color: c}
}

func pointsToVertices(vs []Vector, c sdl.Color) (verts []sdl.Vertex) {
	for _, v := range vs {
		verts = append(verts, pointToVertex(v, c))
	}

	return
}

func drawConnectedPoints(renderer *sdl.Renderer, vs []Vector) {
	firstP := vs[0]
	lastP := vs[0]

	for _, v := range vs {
		renderer.DrawLine(int32(lastP.X), int32(lastP.Y), int32(v.X), int32(v.Y))
		lastP = v
	}

	renderer.DrawLine(int32(firstP.X), int32(firstP.Y), int32(lastP.X), int32(lastP.Y))
}

func drawShape(renderer *sdl.Renderer, shape Shape) {
	renderer.SetDrawColor(shape.Color.R, shape.Color.G, shape.Color.B, shape.Color.A)
	switch shape.ShapeType {
	case ShapeRectFilled:
		if shape.Angle != 0 {
			points := genRectPoints(float32(shape.X), float32(shape.Y), shape.Size, shape.Angle)
			renderer.RenderGeometry(nil, pointsToVertices([]Vector{points[0], points[1], points[3]}, shape.Color), []int32{0, 1, 2})
			renderer.RenderGeometry(nil, pointsToVertices([]Vector{points[1], points[2], points[3]}, shape.Color), []int32{2, 1, 0})
		} else {
			renderer.FillRect(&(sdl.Rect{X: shape.X, Y: shape.Y, W: shape.Size, H: shape.Size}))
		}

	case ShapeRectOutline:
		if shape.Angle != 0 {
			drawConnectedPoints(renderer, genRectPoints(float32(shape.X), float32(shape.Y), shape.Size, shape.Angle))
		} else {
			renderer.DrawRect(&(sdl.Rect{X: shape.X, Y: shape.Y, W: shape.Size, H: shape.Size}))
		}

	case ShapeTriangleFilled:
		renderer.RenderGeometry(nil, pointsToVertices(genTrianglePoints(float32(shape.X), float32(shape.Y), shape.Size, shape.Angle), shape.Color), []int32{0, 1, 2})

	case ShapeTriangleOutline:
		drawConnectedPoints(renderer, genTrianglePoints(float32(shape.X), float32(shape.Y), shape.Size, shape.Angle))
	}
}

func randomColor() sdl.Color {
	return sdl.Color{R: uint8(rand.Uint32()), G: uint8(rand.Uint32()), B: uint8(rand.Uint32()), A: 255}
}

func newShape(renderer *sdl.Renderer, x int32, y int32) {
	ns := Shape{
		ShapeType:      uint8(rand.Uint32() % 4),
		X:              x,
		Y:              y,
		Size:           int32(rand.Uint32()%150 + 25),
		Angle:          (rand.Float32() * 2 * math.Pi),
		Color:          randomColor(),
		MovementVector: Vector{X: rand.Float32() - 0.5, Y: rand.Float32() - 0.5},
	}

	if ns.MovementVector.X <= 0.1 && ns.MovementVector.X >= 0.1 {
		ns.MovementVector.X = 1
	}

	if ns.MovementVector.Y <= 0.1 && ns.MovementVector.Y >= 0.1 {
		ns.MovementVector.Y = 1
	}

	Shapes = append(Shapes, ns)

	drawShape(renderer, ns)
}

func updateShapes(window *sdl.Window) {
	for i := range Shapes {
		checkForScreenCollision(window, &(Shapes[i]))

		Shapes[i].X += int32(Shapes[i].MovementVector.X * 10)
		Shapes[i].Y += int32(Shapes[i].MovementVector.Y * 10)
	}
}

func renderShapes(renderer *sdl.Renderer) {
	renderer.SetDrawColor(0, 0, 0, 255)
	renderer.Clear()

	for _, shape := range Shapes {
		drawShape(renderer, shape)
	}
}

func newTempPoint(x int32, y int32) {
	if len(TempShapes) > 2555 {
		TempShapes = TempShapes[1:]
	}

	var shape uint8
	if rand.Uint32()%2 == 0 {
		shape = ShapeRectOutline
	} else {
		shape = ShapeTriangleOutline
	}

	TempShapes = append(TempShapes, Shape{X: x, Y: y, Color: sdl.Color{R: 25, G: 25, B: 25} /*randomColor()*/, ShapeType: shape, Size: 10, Angle: 0})
}

func updateTempPoints() {
	temp := TempShapes
	TempShapes = []Shape{}
	for _, t := range temp {
		t.Size--
		t.Color.R--
		t.Color.G--
		t.Color.B--

		if t.Size > 0 {
			TempShapes = append(TempShapes, t)
		}
	}
}

func renderTempPoints(renderer *sdl.Renderer) {
	for _, shape := range TempShapes {
		drawShape(renderer, shape)
	}
}

func render(renderer *sdl.Renderer) {
	renderShapes(renderer)
	renderTempPoints(renderer)

	renderer.Present()
}

func checkForScreenCollision(window *sdl.Window, shape *Shape) {
	var points []Vector
	if shape.ShapeType == ShapeRectFilled || shape.ShapeType == ShapeRectOutline {
		points = genRectPoints(float32(shape.X), float32(shape.Y), shape.Size, shape.Angle)
	} else {
		points = genTrianglePoints(float32(shape.X), float32(shape.Y), shape.Size, shape.Angle)
	}

	w, h := window.GetSize()
	for _, p := range points {
		if int32(p.X) > w || p.X < 0 {
			shape.MovementVector.X *= -1.02

			shape.X += int32(shape.MovementVector.X) * 10
			shape.Y += int32(shape.MovementVector.Y) * 10
		} else if int32(p.Y) > h || p.Y < 0 {
			shape.MovementVector.Y *= -1.02

			shape.X += int32(shape.MovementVector.X) * 10
			shape.Y += int32(shape.MovementVector.Y) * 10
		}
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"graphical boredom killer",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN,
	)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	renderer.Present()

	running := true

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONDOWN {
					newShape(renderer, t.X, t.Y)
				}

			case *sdl.MouseMotionEvent:
				newTempPoint(t.X, t.Y)

			case *sdl.QuitEvent:
				running = false
			}
		}

		updateShapes(window)
		updateTempPoints()

		render(renderer)

		sdl.Delay(50)
	}
}
