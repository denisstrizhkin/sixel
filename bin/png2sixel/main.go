package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/denisstrizhkin/sixel"
)

func colors_to_pixels(img image.Image) []sixel.Pixel {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	pixels := make([]sixel.Pixel, 0, w*h)
	for i := range h {
		for j := range w {
			r, g, b, _ := img.At(j, i).RGBA()
			pixels = append(pixels, sixel.Pixel{
				R:       r >> 8,
				G:       g >> 8,
				B:       b >> 8,
				A:       0,
				W:       0,
				Cluster: -1,
			})
		}
	}
	return pixels
}

func save_palette(p []sixel.Pixel) {
	w := 100
	h := 100
	fW := w * len(p)

	upLeft := image.Point{0, 0}
	lowRight := image.Point{fW, h}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	log.Println(p)
	for i := range h {
		for j := range fW {
			p_i := j / w
			c := color.RGBA{
				uint8(p[p_i].R),
				uint8(p[p_i].G),
				uint8(p[p_i].B),
				255,
			}
			img.Set(j, i, c)
		}
	}

	f, _ := os.Create("palette.png")
	png.Encode(f, img)
}

func sixel_encode(img image.Image, w io.Writer) {
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	header := fmt.Sprintf("\x1bPq\"1;1;%d;%d", width, height)
	bw.Write([]byte(header))

	pixels := colors_to_pixels(img)
	palette, clusterMap := sixel.Clusterize(pixels, 256, 100)
	//save_palette(palette)

	for i, p := range palette {
		r := p.R * 100 / 255
		g := p.G * 100 / 255
		b := p.B * 100 / 255
		bw.Write([]byte(fmt.Sprintf("#%d;2;%d;%d;%d", i, r, g, b)))
	}

	for i := range height {
		for j := range width {
			p_id := clusterMap[pixels[i*width+j]]
			c := rune((1 << (i % 6)) + 63)
			bw.Write([]byte(fmt.Sprintf("#%d%c", p_id, c)))
		}
		if i%6 == 5 {
			bw.Write([]byte("-"))
		} else {
			bw.Write([]byte("$"))
		}
	}

	bw.Write([]byte("\x1b\\"))
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalln("no image path provided")
	}

	img_path := os.Args[1]
	img_file, err := os.Open(img_path)
	if err != nil {
		log.Fatalln("error opening file:", err)
	}
	defer img_file.Close()

	img, err := png.Decode(img_file)
	if err != nil {
		log.Fatalln("can't decode png image:", err)
	}

	sixel_encode(img, os.Stdout)
}
