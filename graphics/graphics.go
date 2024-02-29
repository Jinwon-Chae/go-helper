package graphics

import (
	"errors"
	"fmt"

	"github.com/fogleman/gg"
)

type Graphics struct {
	ptr *gg.Context
}

func NewGraphics() *Graphics {
	return &Graphics{}
}

/*
      색상표
___________________
색상 | r | g | b |
-------------------
검정 | 0 | 0 | 0 |
-------------------
하얀 | 1 | 1 | 1 |
-------------------

*/

// @param width int 배경 이미지 가로 픽셀 수
// @param height int 배경 이미지 세로 픽셀 수
// @param red float64 배경 이미지 red 비율
// @param green float64 배경 이미지 green 비율
// @param blue float64 배경 이미지 blue 비율
func (g *Graphics) Open(width, height int, red, green, blue float64) (err error) {
	if width <= 0 || height <= 0 {
		return errors.New("graphic open fail: width or height is invalid")
	}

	g.ptr = gg.NewContext(width, height)

	// red, green, blue는 내부에서 uint로 타입 변환을 하기 때문에 유효 체크 별도 안함
	g.ptr.SetRGB(red, green, blue)
	g.ptr.Clear()

	return
}

func (g *Graphics) Close() {
	g.ptr = nil
}

// @param name string 저장 이미지 이름
func (g *Graphics) SavePNG(name string) (err error) {
	return g.ptr.SavePNG(fmt.Sprintf("%s.png", name))
}
