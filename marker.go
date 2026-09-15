// xsscolorpicker/marker.go

package colorpicker

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Цвета маркеров берутся из currentStyle (см. style.go). Переменные
// markerFillColor/markerStrokeColor удалены — раньше они были глобальным
// хардкодом и не подчинялись теме.

type marker interface {
	fyne.CanvasObject

	position() fyne.Position
	setPosition(fyne.Position)
	object() fyne.CanvasObject
}

func setPositionX(m marker, x float32) {
	m.setPosition(fyne.NewPos(x, m.position().Y))
}

func setPositionY(m marker, y float32) {
	m.setPosition(fyne.NewPos(m.position().X, y))
}

type defaultMarker struct {
	*canvas.Circle
	center fyne.Position
	radius float32
}

func newDefaultMarker(radius float32) marker {
	marker := &defaultMarker{
		Circle: &canvas.Circle{
			FillColor:   currentStyle.MarkerFill,
			StrokeColor: currentStyle.MarkerStroke,
			StrokeWidth: 1,
		},
		radius: radius,
	}
	marker.setPosition(fyne.NewPos(0, 0))
	return marker
}

func (m *defaultMarker) position() fyne.Position {
	return m.center
}

func (m *defaultMarker) setPosition(p fyne.Position) {
	m.center = p
	topLeft := fyne.NewPos(p.X-m.radius, p.Y-m.radius)
	size := fyne.NewSize(m.radius*2, m.radius*2)
	m.Circle.Move(topLeft)
	m.Circle.Resize(size)
}

func (m *defaultMarker) object() fyne.CanvasObject {
	return m.Circle
}

type barMarker interface {
	marker

	setPositionFromValue(v float32)
	calcValueFromPosition(p fyne.Position) float32
}

type defaultBarMarker struct {
	marker
}

func newDefaultBarMarker(barWidth float32) barMarker {
	m := newDefaultMarker(barWidth / 2)
	return &defaultBarMarker{marker: m}
}

func (m *defaultBarMarker) setPositionFromValue(v float32) {
	panic("not implemented")
}

func (m *defaultBarMarker) calcValueFromPosition(p fyne.Position) float32 {
	panic("not implemented")
}

type circleBarMarker struct {
	*defaultMarker
	cx, cy float32
}

func newCircleBarMarker(w, h float32, barWidth float32) *circleBarMarker {
	fw := float64(w)
	fh := float64(h)
	fr := barWidth / 2
	marker := &circleBarMarker{
		defaultMarker: newDefaultMarker(fr).(*defaultMarker),
		cx:            w / 2,
		cy:            h / 2,
	}
	markerCenter := fyne.NewPos(float32(math.Round(fw-float64(fr))), float32(math.Round(fh/2)))
	marker.defaultMarker.setPosition(markerCenter)
	return marker
}

func (m *circleBarMarker) setPosition(p fyne.Position) {
	v := newVectorFromPoints(float64(m.cx), float64(m.cy), float64(p.X), float64(p.Y))
	nv := v.normalize()
	center := newVector(float64(m.cx), float64(m.cy))
	markerCenter := center.add(nv.multiply(float64(m.cx - m.radius))).toPosition()
	m.defaultMarker.setPosition(markerCenter)
}

func (m *circleBarMarker) setPositionFromValue(v float32) {
	hue := v
	rad := float64(-2 * math.Pi * hue)
	center := newVector(float64(m.cx), float64(m.cy))
	dir := newVector(1, 0).rotate(rad).multiply(float64(m.cx - m.radius))
	markerCenter := center.add(dir).toPosition()
	m.defaultMarker.setPosition(markerCenter)
}

func (m *circleBarMarker) calcValueFromPosition(p fyne.Position) float32 {
	v := newVectorFromPoints(float64(m.cx), float64(m.cy), float64(p.X), float64(p.Y))
	baseV := newVector(1, 0)
	rad := math.Acos(baseV.dot(v) / (v.norm() * baseV.norm()))
	if float64(p.Y-m.cy) >= 0 {
		rad = math.Pi*2 - rad
	}
	hue := rad / (math.Pi * 2)
	return float32(hue)
}
