package graphics

type ImageArea struct {
	Width  int // 가로 픽셀 수
	Height int // 세로 픽셀 수
}

func (i *ImageArea) isInvalid() bool {
	if i.Width <= 0 || i.Height <= 0 {
		return false
	}

	return true
}

/*   색상표
___________________
색상 | r | g | b |
-------------------
검정 | 0 | 0 | 0 |
-------------------
하얀 | 1 | 1 | 1 |
-------------------
빨강 | 1 | 0 | 0 |
-------------------
초록 | 0 | 1 | 0 |
-------------------
파랑 | 0 | 0 | 1 |
-------------------
노랑 | 1 | 1 | 0 |
-------------------*/

type Color struct {
	Red   float64 // RGB 중 red 비율
	Green float64 // RGB 중 green 비율
	Blue  float64 // RGB 중 red 비율
}

type Font string

const (
	Arial           Font = "font/Arial.ttf"
	NanumGothic     Font = "font/NanumGothic.ttf"
	NanumGothicBold Font = "font/NanumGothicBold.ttf"
)

type Position struct {
	X float64 // 문자 시작 x 좌표
	Y float64 // 문자 시작 y 좌표
}

type Anchor struct { // 앵커 포인트
	X float64
	Y float64
}

type TextConfig struct {
	Color    Color    // 문자열 RGB 비중
	Font     Font     // 폰트
	Size     int      // 문자 사이즈
	Position Position // 문자 시작 좌표
	Anchor   Anchor   // 앵커포인트
}
