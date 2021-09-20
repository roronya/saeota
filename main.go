package main

import (
	"flag"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
)

const (
	WIDTH, HEIGHT = 1280, 720
	FIGURE_POINT_MIN_X = 240
	FIGURE_POINT_MIN_Y = 10
	FIGURE_POINT_MAX_X = 1040
	FIGURE_POINT_MAX_Y = 466
	FONT_SIZE = 26
	// コメントの場所は何度か試して適当に決めた
	R_COMMENT_POINT_X = 250
	R_COMMENT_POINT_Y = 530
	R_COMMENT2_POINT_Y = R_COMMENT_POINT_Y + FONT_SIZE + 5 // 5は行間のサイズ
	L_COMMENT_POINT_X = 250
	L_COMMENT_POINT_Y = 645
	L_COMMENT2_POINT_Y = L_COMMENT_POINT_Y + FONT_SIZE + 5
)

var l,r,f string

func init() {
	flag.StringVar(&l, "l", "", "左のセリフ")
	flag.StringVar(&r, "r", "", "右のセリフ")
	flag.StringVar(&f, "f", "", "解説対象の画像パス")
	flag.Parse()
	if f == "" {
		log.Fatal("解説対象の画像パスは必ず指定してください")
	}
}

func main () {
	file, err := os.Open("assets/lr.png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	template, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	file,err = os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	figure, err := png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// 新しく生成する画像を初期化
	dst := image.NewRGBA(
		image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: WIDTH, Y: HEIGHT},
	})
	// templateの画像を上に乗せる
	draw.Draw(
		dst,
		// 第2引数は下地になる画像のどの範囲に上に乗せるかを指定する。templateで完全に覆うので(0,0)と(1280,720)を指定する
		image.Rectangle{Min: image.Point{}, Max: image.Point{X: WIDTH, Y: HEIGHT}},
		template,
		//　第4引数は上に乗せる画像を指定したpointから切り取る。切り取る必要はないので(0, 0)を指定する
		image.Point{},
		draw.Src,
	)
	// 解説対象の画像をtemplateの上に乗せる
	draw.Draw(
		dst,
		// 解説対象の画像は800*450なので左からは(1280-800)/2=240pxの場所に置く。画面上部のマージンはなんとなく10pxくらい
		// 右下は左隅から右に240+800=1040pxで、縦は16+450=466pxを指定する
		image.Rectangle{
			Min: image.Point{X: FIGURE_POINT_MIN_X, Y: FIGURE_POINT_MIN_Y},
			Max: image.Point{X: FIGURE_POINT_MAX_X, Y: FIGURE_POINT_MAX_Y}},
		figure,
		image.Point{}, // 切り取らないから(0,0)を指定する
		draw.Src,
	)

	ttf, err := os.Open("assets/ipaexg.ttf")
	if err != nil {
		log.Fatal(err)
	}
	ftBin, err := io.ReadAll(ttf)
	if err != nil {
		log.Fatal(err)
	}
	ft, err := truetype.Parse(ftBin)
	if err != nil {
		log.Fatal(err)
	}

	drawComment(r, ft, R_COMMENT_POINT_X, R_COMMENT_POINT_Y, dst)
	drawComment(l, ft, L_COMMENT_POINT_X, L_COMMENT_POINT_Y, dst)

	png.Encode(os.Stdout, dst)
}

func drawComment(text string, ft *truetype.Font, x int, y int, dst *image.RGBA) {
	opt := truetype.Options{Size: FONT_SIZE}
	face := truetype.NewFace(ft, &opt)
	dot := fixed.Point26_6{X: fixed.Int26_6(x*64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst: dst,
		Src: image.NewUniform(color.Black),
		Face: face,
		Dot: dot,
	}
	d.DrawString(text)
}