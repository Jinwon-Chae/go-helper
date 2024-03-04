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

// @param imageArea ImageArea(strct) 배경 이미지 가로 세로 픽셀 수
// @param color Color(struct) 배경 이미지 색
func (g *Graphics) Open(imageArea ImageArea, color Color) (err error) {
	if !imageArea.isInvalid() {
		return errors.New("graphic open fail: width or height is invalid")
	}

	g.ptr = gg.NewContext(imageArea.Width, imageArea.Height)

	// red, green, blue는 내부에서 uint로 타입 변환을 하기 때문에 유효 체크 별도 안함
	g.ptr.SetRGB(color.Red, color.Green, color.Blue)
	g.ptr.Clear()

	return
}

func (g *Graphics) Close() {
	g.ptr = nil
}

// @param config TextConfig(strct) 텍스트 설정 값
// @param text string 텍스트 내용
func (g *Graphics) WriteText(config TextConfig, text string) (err error) {
	g.ptr.SetRGB(config.Color.Red, config.Color.Green, config.Color.Blue)
	if err = g.ptr.LoadFontFace(string(config.Font), float64(config.Size)); err != nil {
		return errors.New("write text fail: load font face fail -> " + err.Error())
	}

	g.ptr.DrawString(text, config.Position.X, config.Position.Y)

	return
}

// @param name string 저장 이미지 이름(확장자 제외)
func (g *Graphics) SavePNG(name string) (err error) {
	return g.ptr.SavePNG(fmt.Sprintf("%s.png", name))
}
