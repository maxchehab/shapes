package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/png"
	"os"

	"github.com/andybons/gogif"
	noise "github.com/ojrac/opensimplex-go"
)

type Pixel struct {
	color color.RGBA
	x     int
	y     int
}

const WIDTH = 1080
const HEIGHT = 1920
const SEED = 7150434019954586410
const FRAME_COUNT = 60

func main() {
	var frames = make([][WIDTH][HEIGHT]float64, 0)

	var lastFrame = generateGrid()
	fmt.Println("Created initial frame (noise)")

	frames = append(frames, lastFrame)

	for i := 0; i < FRAME_COUNT-1; i++ {
		lastFrame = createNextFrame(lastFrame, i)
		frames = append(frames, lastFrame)
	}

	saveGif(frames)
}

func createNextFrame(frame [WIDTH][HEIGHT]float64, index int) [WIDTH][HEIGHT]float64 {
	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			var noise = frame[x][y]

			if index > 30 {
				frame[x][y] = noise + 1
			} else {
				frame[x][y] = noise - 1
			}
		}
	}

	return frame
}

func generateGrid() (frame [WIDTH][HEIGHT]float64) {
	noise := noise.New(SEED)

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			var xFloat = float64(x) / float64(WIDTH)
			var yFloat = float64(y) / float64(HEIGHT)

			frame[x][y] = noise.Eval2(xFloat, yFloat)
		}
	}

	return frame
}

func imageFromGrid(grid [WIDTH][HEIGHT]Pixel) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	for _, row := range grid {
		for _, pixel := range row {
			img.Set(pixel.x, pixel.y, pixel.color)
		}
	}

	return img
}

func saveGif(frames [][WIDTH][HEIGHT]float64) {
	outGif := &gif.GIF{}

	for i, frame := range frames {
		fmt.Printf("Generating frame %v\n", i)

		var img = image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

		for x, row := range frame {
			for y, noise := range row {
				c := uint8(noise * 255)
				img.Set(x, y, color.RGBA{c, c, c, 255})
			}
		}

		bounds := img.Bounds()
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer := gogif.MedianCutQuantizer{NumColor: 255}
		quantizer.Quantize(palettedImage, bounds, img, image.ZP)

		outGif.Image = append(outGif.Image, palettedImage)
		outGif.Delay = append(outGif.Delay, 0)
	}

	// save to out.gif
	f, _ := os.OpenFile("map.gif", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	gif.EncodeAll(f, outGif)
}

func savePng(img image.Image) {
	f, _ := os.OpenFile("map.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
