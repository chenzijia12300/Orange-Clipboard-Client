package ui

import (
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

var (
	Blue = color.RGBA{
		R: 173,
		G: 216,
		B: 230,
		A: 0,
	}
	Green = color.RGBA{
		R: 152,
		G: 255,
		B: 152,
		A: 100,
	}
	DefaultColor = Blue
)

var (
	DefaultFontSize float32 = 12.0
)

func createText(text string, color color.Color, fontSize float32) *canvas.Text {
	textWidget := canvas.NewText(text, color)
	textWidget.TextSize = fontSize
	return textWidget
}

func CreateDefaultText(text string) *canvas.Text {
	return createText(text, color.Black, DefaultFontSize)
}
