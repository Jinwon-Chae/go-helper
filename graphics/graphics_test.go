package graphics

import (
	"testing"

	"github.com/fogleman/gg"
)

func TestDrawBackground(t *testing.T) {

}

func TestDrawText(t *testing.T) {

}

func TestGGtest(t *testing.T) {
	const S = 1024
	dc := gg.NewContext(200, 20)
	dc.SetRGB(0, 0, 0)
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace("Arial.ttf", 20); err != nil {
		panic(err)
	}
	dc.DrawString("holllll", 0, 20)
	// dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	dc.SavePNG("out.png")
}
