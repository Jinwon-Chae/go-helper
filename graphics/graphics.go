package graphics

import (
	"errors"
	"fmt"
	"image"
	"math"

	"github.com/fogleman/gg"
)

type Graphics struct {
	ptr *gg.Context
}

func NewGraphics() *Graphics {
	return &Graphics{}
}

// 최초 배경 이미지 생성 및 설정
// @param imageArea ImageArea(strct): 배경 이미지 가로 세로 픽셀 수
// @param color Color(struct): 배경 이미지 색
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

// 인스턴스 명시적 소멸
func (g *Graphics) Close() {
	g.ptr = nil
}

// 텍스트 설정
// @param config TextConfig(strct): 텍스트 설정 값
// @param text string: 텍스트 내용
func (g *Graphics) WriteText(config TextConfig, text string) (err error) {
	g.ptr.SetRGB(config.Color.Red, config.Color.Green, config.Color.Blue)
	if err = g.ptr.LoadFontFace(string(config.Font), float64(config.Size)); err != nil {
		return errors.New("write text fail: load font face fail -> " + err.Error())
	}

	g.ptr.DrawString(text, config.Position.X, config.Position.Y)

	return
}

// 기존 설정 초기화
func (g *Graphics) Clear() {
	g.ptr.Clear()
}

// 새롭게 배경 이미지 색 설정
func (g *Graphics) SetBackgroundColor(color Color) {
	g.ptr.SetRGB(color.Red, color.Green, color.Blue)
	g.ptr.Clear()
}

// 이미지 저장
// @param name string: 저장 이미지 이름(확장자 제외)
func (g *Graphics) SavePNG(name string) (err error) {
	return g.ptr.SavePNG(fmt.Sprintf("%s.png", name))
}

// 이미지 수평 병합
// @param name string: 병합된 파일 이름(확장자 제외)
// @param paths ...string: 수평 평합 대상 이미지 경로들
func HorizontalConcatPNG(name string, paths ...string) (err error) {
	if len(paths) < 1 {
		return errors.New("path is empty")
	}

	var width int
	var height int
	var images []image.Image
	for _, path := range paths {
		i, err := gg.LoadPNG(path)
		if err != nil {
			return errors.New("png load fail: " + path)
		}

		s := i.Bounds().Size()

		width += s.X
		height = int(math.Max(float64(height), float64(s.Y)))

		images = append(images, i)
	}

	var beforeX int
	dc := gg.NewContext(width, height)
	for _, image := range images {
		dc.DrawImage(image, beforeX, 0)
		beforeX = image.Bounds().Size().X
	}

	if err = dc.SavePNG(name); err != nil {
		return errors.New("save concat png fail: " + err.Error())
	}

	return
}

// 이미지 수직 병합
// @param name string: 병합된 파일 이름(확장자 제외)
// @param paths ...string: 수직 평합 대상 이미지 경로들
func VerticalConcatPNG(name string, paths ...string) (err error) {
	if len(paths) < 1 {
		return errors.New("path is empty")
	}

	var width int
	var height int
	var images []image.Image
	for _, path := range paths {
		i, err := gg.LoadPNG(path)
		if err != nil {
			return errors.New("png load fail: " + path)
		}

		s := i.Bounds().Size()

		width = int(math.Max(float64(width), float64(s.X)))
		height += s.Y

		images = append(images, i)
	}

	var beforeY int
	dc := gg.NewContext(width, height)
	for _, image := range images {
		dc.DrawImage(image, 0, beforeY)
		beforeY = image.Bounds().Size().Y
	}

	if err = dc.SavePNG(name); err != nil {
		return errors.New("save concat png fail: " + err.Error())
	}

	return
}
