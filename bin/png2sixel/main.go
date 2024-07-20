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

func get_colors_slice(img image.Image) []color.RGBA {
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	colors := make([]color.RGBA, w*h)

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			colors[i*w+j].R = uint8(r)
			colors[i*w+j].G = uint8(g)
			colors[i*w+j].B = uint8(b)
			colors[i*w+j].A = uint8(a)
		}
	}

	return colors
}

func sixel_encode(img image.Image, w io.Writer) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	header := fmt.Sprintf("\x1bPq1;1;%d;%d", width, height)
	w.Write([]byte(header))

	colors := get_colors_slice(img)
	km := sixel.NewKMeans(256)
	km.Clusterize(colors)

	w.Write([]byte("\x1b\\"))
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

	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	sixel_encode(img, bw)
}
