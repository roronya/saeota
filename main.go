package main

import (
	"bytes"
	_ "embed"
	"flag"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
)

// 画像サイズ
const (
	HEIGHT = 720
	WIDTH  = 1280
)

// 解説対象の画像の表示位置
const (
	FigurePointMinX = 240
	FigurePointMinY = 10
	FigurePointMaxX = 1040
	FigurePointMaxY = 466
)

// セリフを描画するフォントや位置 何度か試して適当に決めた
const (
	FontSize        = 26
	RCommentPointX  = 250
	RCommentPointY  = 530
	RComment2PointY = RCommentPointY + FontSize + 5 // 5は行間のサイズ
	LCommentPointX  = 250
	LCommentPointY  = 645
	LComment2PointY = LCommentPointY + FontSize + 5
)

/* 変数は全てinitで初期化してmainを見通し良くする */
var l, l2, r, r2 string // セリフ
var f string            // 解説対象の画像のファイルパス
var figure image.Image  // 解説対象の画像

// テンプレート画像の読み込み用
//go:embed assets/lr.png
var template_lr []byte

//go:embed assets/l.png
var template_l []byte

//go:embed assets/r.png
var template_r []byte

//go:embed assets/nocomment.png
var template_nocomment []byte

var template image.Image // テンプレートの画像

//go:embed assets/ipaexg00401/ipaexg.ttf
var ftBin []byte // フォントの読み込み用

var ft *truetype.Font // セリフを描画するときに使うフォント

func init() {
	// コマンドライン引数
	flag.StringVar(&l, "l", "", "左のセリフ")
	flag.StringVar(&l2, "l2", "", "二行目の左のセリフ")
	flag.StringVar(&r, "r", "", "右のセリフ")
	flag.StringVar(&r2, "r2", "", "二行目の右のセリフ")
	flag.StringVar(&f, "f", "", "解説対象の画像パス")
	flag.Parse()
	if f == "" {
		log.Fatal("解説対象の画像パスは必ず指定してください")
	}

	// 解説対象の画像の取得
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	figure, err = png.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	// セリフのありなしを見てどのテンプレートを使うか決めてオープンする
	t := template_nocomment
	if l != "" && r != "" {
		t = template_lr
	}
	if l == "" && r != "" {
		t = template_r
	}
	if l != "" && r == "" {
		t = template_l
	}
	template, err = png.Decode(bytes.NewReader(t))
	if err != nil {
		log.Fatal(err)
	}

	// フォントを開いて用意する
	ft, err = truetype.Parse(ftBin)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
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
			Min: image.Point{X: FigurePointMinX, Y: FigurePointMinY},
			Max: image.Point{X: FigurePointMaxX, Y: FigurePointMaxY}},
		figure,
		image.Point{}, // 切り取らないから(0,0)を指定する
		draw.Src,
	)

	if r != "" {
		drawComment(r, ft, RCommentPointX, RCommentPointY, dst)
	}
	if r2 != "" {
		drawComment(r2, ft, RCommentPointX, RComment2PointY, dst)
	}
	if l != "" {
		drawComment(l, ft, LCommentPointX, LCommentPointY, dst)
	}
	if l2 != "" {
		drawComment(l2, ft, LCommentPointX, LComment2PointY, dst)
	}

	png.Encode(os.Stdout, dst)
}

func drawComment(text string, ft *truetype.Font, x int, y int, dst *image.RGBA) {
	opt := truetype.Options{Size: FontSize}
	face := truetype.NewFace(ft, &opt)
	dot := fixed.Point26_6{X: fixed.Int26_6(x * 64), Y: fixed.Int26_6(y * 64)}
	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  dot,
	}
	d.DrawString(text)
}
