package main

import (
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main () {
	const (
		width, height = 1280, 720
	)
	f := "img/test.png"

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
		Max: image.Point{X: width, Y: height},
	})
	// templateの画像を上に乗せる
	draw.Draw(
		dst,
		// 第2引数は下地になる画像のどの範囲に上に乗せるかを指定する。templateで完全に覆うので(0,0)と(1280,720)を指定する
		image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}},
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
		image.Rectangle{Min: image.Point{X: 240, Y: 10}, Max: image.Point{X: 1040, Y: 466}},
		figure,
		image.Point{}, // 切り取らないから(0,0)を指定する
		draw.Src,
	)
	png.Encode(os.Stdout, dst)
}