package graphics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDrawText(t *testing.T) {
	g := NewGraphics()
	defer g.Close()
	assert.NotNil(t, g)

	i := ImageArea{Width: 320, Height: 32} // 내민식
	c := Color{Red: 0, Green: 0, Blue: 0}  // 검정
	err := g.Open(i, c)
	assert.Nil(t, err)

	c = Color{Red: 1, Green: 1, Blue: 1}
	p := Position{X: 160, Y: 16}
	config := TextConfig{Color: c, Font: NanumGothicBold, Size: 20, Position: p}
	text1 := "교통"
	err = g.WriteText(config, text1)
	assert.Nil(t, err)

	c = Color{Red: 1, Green: 1, Blue: 0}
	p = Position{X: 0, Y: 0}
	config = TextConfig{Color: c, Font: NanumGothicBold, Size: 20, Position: p}
	text2 := "고통"

	err = g.WriteText(config, text2)
	assert.Nil(t, err)

	file := "out2"
	err = g.SavePNG(file)
	assert.Nil(t, err)
}

func TestConcat(t *testing.T) {
	g := NewGraphics()
	defer g.Close()
	assert.NotNil(t, g)

	i := ImageArea{Width: 160, Height: 32}
	c := Color{Red: 0, Green: 0, Blue: 0} // 검정
	err := g.Open(i, c)
	assert.Nil(t, err)

	c = Color{Red: 1, Green: 1, Blue: 1}
	p := Position{X: 0, Y: 20}
	config := TextConfig{Color: c, Font: NanumGothicBold, Size: 20, Position: p}
	text1 := "교통"
	err = g.WriteText(config, text1)
	assert.Nil(t, err)

	file := "out1"
	err = g.SavePNG(file)
	assert.Nil(t, err)

	g.Clear()
	g.SetBackgroundColor(Color{Red: 0, Green: 0, Blue: 0})
	c = Color{Red: 1, Green: 1, Blue: 0}
	p = Position{X: 40, Y: 20}
	config = TextConfig{Color: c, Font: NanumGothicBold, Size: 20, Position: p}
	text2 := "고통"
	err = g.WriteText(config, text2)
	assert.Nil(t, err)

	file = "out2"
	err = g.SavePNG(file)
	assert.Nil(t, err)

	err = HorizontalConcatPNG("text1.png", "out1.png", "out2.png")
	assert.Nil(t, err)

	err = VerticalConcatPNG("test.png", "out1.png", "out2.png")
	assert.Nil(t, err)
}
