package strtoimg

import (
	"image"
	"image/color"
	"image/draw"
)

// color.Black or color.White
func drawBackground(width, height int, color color.Gray16) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)

	return img
}

// func textToImage(text string) image.Image {
// 	// 이미지 크기 설정
// 	width := 200
// 	height := 50

// 	// 흰 배경의 이미지 생성
// 	img := image.NewRGBA(image.Rect(0, 0, width, height))
// 	white := color.RGBA{255, 255, 255, 255}
// 	draw.Draw(img, img.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

// 	// // 이미지에 문자열 그리기
// 	// // 폰트 및 위치 등을 조정하여 원하는 스타일로 표시 가능
// 	// basicFont := &image.Uniform{color.Black}
// 	// drawText(img, 10, 20, text, basicFont)

// 	return img
// }
