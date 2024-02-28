package strtoimg

import (
	"image/color"
	"image/png"
	"os"
	"testing"
)

func TestDrawBackground(t *testing.T) {
	img := drawBackground(200, 50, color.Black)

	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// PNG 형식으로 이미지 저장
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
