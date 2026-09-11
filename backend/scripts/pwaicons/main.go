package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func icon(size int) image.Image {
	wood := color.NRGBA{R: 0x1a, G: 0x14, B: 0x10, A: 255}
	paper := color.NRGBA{R: 0xf7, G: 0xf1, B: 0xe8, A: 255}
	chili := color.NRGBA{R: 0xb4, G: 0x23, B: 0x18, A: 255}
	dst := imaging.New(size, size, wood)
	insetX := size * 10 / 64
	insetY := size * 12 / 64
	sheetW := size * 44 / 64
	sheetH := size * 40 / 64
	sheet := imaging.New(sheetW, sheetH, paper)
	dst = imaging.Overlay(dst, sheet, image.Pt(insetX, insetY), 1)
	lineW := max(2, size*24/640)
	for i, row := range []int{22, 32, 42} {
		y := size * row / 64
		w := []int{28, 20, 16}[i] * size / 64
		x0 := size * 18 / 64
		for x := x0; x < x0+w; x++ {
			for t := 0; t < lineW; t++ {
				dst.SetNRGBA(x, y+t, chili)
			}
		}
	}
	return dst
}

func writePNG(path string, img image.Image) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func main() {
	root := filepath.Join("..", "..", "frontend", "public")
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	for name, img := range map[string]image.Image{
		"pwa-192.png":          icon(192),
		"pwa-512.png":          icon(512),
		"apple-touch-icon.png": icon(180),
	} {
		if err := writePNG(filepath.Join(root, name), img); err != nil {
			panic(err)
		}
	}
}
