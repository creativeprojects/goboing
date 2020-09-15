package main

import (
	"log"

	"github.com/hajimehoshi/ebiten"
)

// XType represents the type of the X coordinate (centre, left or right)
type XType int

// XType
const (
	XCentre XType = iota
	XLeft
	XRight
)

// YType represents the type of the Y coordinate (centre, top or bottom)
type YType int

// YType
const (
	YCentre YType = iota
	YTop
	YBottom
)

// Sprite manages sprite movement and animation
type Sprite struct {
	xType     XType
	yType     YType
	frame     int
	x         float64
	y         float64
	image     *ebiten.Image
	animation []*ebiten.Image
	sequence  []int
	rate      int
	op        *ebiten.DrawImageOptions
}

// NewSprite creates a new Sprite with default coordinate type
func NewSprite(xType XType, yType YType) *Sprite {
	return &Sprite{
		xType: xType,
		yType: yType,
		op:    &ebiten.DrawImageOptions{},
	}
}

// SetImage sets sprite image
func (s *Sprite) SetImage(image *ebiten.Image) *Sprite {
	s.image = image
	return s
}

// Update animation (if needed)
func (s *Sprite) Update() {
	s.frame++
}

// Draw the current image to the screen. If no image has been set, it does nothing
func (s *Sprite) Draw(screen *ebiten.Image) {
	if s.image == nil {
		log.Println("Sprite.Draw: no image to draw")
		return
	}
	width, height := s.image.Size()
	s.op.GeoM.Reset()
	s.op.GeoM.Translate(s.xleft(float64(width)), s.ytop(float64(height)))
	screen.DrawImage(s.image, s.op)
}

// Start (or restart) an animation
func (s *Sprite) Start() *Sprite {
	s.frame = 0
	return s
}

// Move to relative coordinates (adds coordinates to the current position)
func (s *Sprite) Move(x, y float64) *Sprite {
	s.x += x
	s.y += y
	return s
}

// MoveTo the new coordinates using the default coordinates type defined at instantiation
func (s *Sprite) MoveTo(x, y float64) *Sprite {
	s.x = x
	s.y = y
	return s
}

// MoveToType moves to the new coordinates using the specified coordinate types
func (s *Sprite) MoveToType(x, y float64, xType XType, yType YType) *Sprite {
	if s.xType == xType {
		s.x = x
	}
	if s.yType == yType {
		s.y = y
	}
	panic("Unfinished")
	return s
}

// X returns x position. If not image is available to calculate width, it returns -1
func (s *Sprite) X(xType XType) float64 {
	if s.image == nil {
		return -1
	}
	width, _ := s.image.Size()
	switch xType {
	case XCentre:
		return s.xcentre(float64(width))
	case XRight:
		return s.xright(float64(width))
	default:
		return s.xleft(float64(width))
	}
}

// Y returns y position. If not image is available to calculate height, it returns -1
func (s *Sprite) Y(yType YType) float64 {
	if s.image == nil {
		return -1
	}
	_, height := s.image.Size()
	switch yType {
	case YCentre:
		return s.ycentre(float64(height))
	case YBottom:
		return s.ybottom(float64(height))
	default:
		return s.ytop(float64(height))
	}
}

func (s *Sprite) xleft(width float64) float64 {
	switch s.xType {
	case XCentre:
		return s.x - (width / 2)
	case XRight:
		return s.x - width
	default:
		return s.x
	}
}

func (s *Sprite) xcentre(width float64) float64 {
	switch s.xType {
	case XCentre:
		return s.x
	case XRight:
		return s.x - (width / 2)
	default:
		return s.x + (width / 2)
	}
}

func (s *Sprite) xright(width float64) float64 {
	switch s.xType {
	case XCentre:
		return s.x + (width / 2)
	case XRight:
		return s.x
	default:
		return s.x + width
	}
}

func (s *Sprite) ytop(height float64) float64 {
	switch s.yType {
	case YCentre:
		return s.y - (height / 2)
	case YBottom:
		return s.y - height
	default:
		return s.y
	}
}

func (s *Sprite) ycentre(height float64) float64 {
	switch s.yType {
	case YCentre:
		return s.y
	case YBottom:
		return s.y - (height / 2)
	default:
		return s.y + (height / 2)
	}
}

func (s *Sprite) ybottom(height float64) float64 {
	switch s.yType {
	case YCentre:
		return s.y + (height / 2)
	case YBottom:
		return s.y
	default:
		return s.y + height
	}
}
