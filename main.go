package main

import (
	"fmt"
	"github.com/andlabs/ui"
	"image"
	"image/draw"
	"time"
)

var w ui.Window

func initUI() {
	button := ui.NewButton("Load")
	button.OnClicked(func() {
		ui.OpenFile(w, func(filename string) {
			go renderPreview(filename)
		})
	})

	stack := ui.NewVerticalStack(button)

	w = ui.NewWindow("Video", 800, 600, stack)
	w.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	w.Show()
}

type imageArea image.RGBA

func (img *imageArea) Paint(rect image.Rectangle) *image.RGBA {
	return (*image.RGBA)(img).SubImage(rect).(*image.RGBA)
}
func (*imageArea) Mouse(ui.MouseEvent)            { return }
func (*imageArea) Key(ui.KeyEvent) (handled bool) { return }

var w2 ui.Window

func showPreview(img *image.RGBA) {
	w2 = ui.NewWindow("Stuff", img.Bounds().Dx(), img.Bounds().Dy(), ui.NewArea(img.Bounds().Dx(), img.Bounds().Dy(), (*imageArea)(img)))
	w2.OnClosing(func() bool {
		ui.Stop()
		return true
	})
	w2.Show()
}

func renderPreview(filename string) {
	data, err := Probe(filename)
	if err != nil {
		panic(err)
	}

	d := data.Format.Duration()

	const PreviewSize = 128

	canvas := image.NewRGBA(image.Rect(0, 0, int(PreviewSize*d/time.Second), PreviewSize))

	for x := canvas.Bounds().Min.X; x < canvas.Bounds().Max.X; {
		img, err := Frame(filename, time.Duration(x)*time.Second/PreviewSize, "-vf", fmt.Sprintf("scale=-1:%d", PreviewSize))
		if err != nil {
			break
		}

		draw.Draw(canvas, img.Bounds().Add(image.Point{x, 0}), img, image.ZP, draw.Src)
		x += img.Bounds().Dx()
	}

	go ui.Do(func() {
		showPreview(canvas)
	})
}

func main() {
	go ui.Do(initUI)

	err := ui.Go()
	if err != nil {
		panic(err)
	}
}
