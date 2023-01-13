package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math/rand"
	"os"
)

type Dither struct {
	sourceImage   *image.Gray // pointer to the source image in grayscale
	width, height int         // dimensions of the source image
	newImage      draw.Image  // resulting dithered image
	outputFormat  string      // format of the output image (e.g. "png", "jpeg")
	outputName    string      // name of the output image file
	outputQuality int         // quality of the output image
	threshold     int         // threshold value used in dithering algorithm
}

func (d *Dither) SaveFile() {
	outputFileName := d.outputName + "." + d.outputFormat
	outputFile, _ := os.Create(outputFileName)
	defer outputFile.Close()
	jpeg.Encode(outputFile, d.newImage, &jpeg.Options{Quality: d.outputQuality})
}

// Default dithering
func (d *Dither) OrderedDither4() {
	dots := [][]int{{64, 128}, {192, 0}}
	for row := 0; row < d.height; row++ {
		for col := 0; col < d.width; col++ {
			dotrow := 1
			if row%2 == 0 {
				dotrow = 0
			}
			dotcol := 1
			if col%2 == 0 {
				dotcol = 0
			}
			px := d.getPixel(col, row)
			if px > dots[dotrow][dotcol] {
				d.newImage.Set(col, row, color.White)
			} else {
				d.newImage.Set(col, row, color.Black)
			}
		}
	}
}

func (d *Dither) OrderedDither9() {
	dots := [][]int{{0, 196, 84}, {168, 140, 56}, {112, 28, 224}}
	for row := 0; row < d.height; row++ {
		for col := 0; col < d.width; col++ {
			dotrow := 0
			if row%3 == 0 {
				dotrow = 2
			} else if row%2 == 0 {
				dotrow = 1
			}
			dotcol := 0
			if col%3 == 0 {
				dotcol = 2
			} else if col%2 == 0 {
				dotcol = 1
			}
			px := d.getPixel(col, row)
			if px > dots[dotrow][dotcol] {
				d.newImage.Set(col, row, color.White)
			} else {
				d.newImage.Set(col, row, color.Black)
			}
		}
	}
}

func (d *Dither) ThresholdDither() {
	if d.threshold == 0 {
		pxList := d.sourceImage.Pix
		d.threshold = 0
		for i := 0; i < len(pxList); i++ {
			d.threshold += int(pxList[i])
		}
		d.threshold = d.threshold / len(pxList)
	}
	for row := 0; row < d.height; row++ {
		for col := 0; col < d.width; col++ {
			px := d.getPixel(col, row)
			if px > d.threshold {
				d.newImage.Set(col, row, color.White)
			} else {
				d.newImage.Set(col, row, color.Black)
			}
		}
	}
}

func (d *Dither) RandomDither() {
	for row := 0; row < d.height; row++ {
		for col := 0; col < d.width; col++ {
			px := d.getPixel(col, row)
			rand := rand.Intn(255)
			if px > rand {
				d.newImage.Set(col, row, color.White)
			} else {
				d.newImage.Set(col, row, color.Black)
			}
		}
	}
}

func (d *Dither) getPixel(x int, y int) int {
	if x > d.width || y > d.height {
		return 0
	}
	r, g, b, _ := d.sourceImage.At(x, y).RGBA()
	res := uint8((r + g + b) / 3)
	return int(res)
}
