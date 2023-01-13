package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"log"
	"os"
)

func main() {
	var method string
	var outputFormat string
	var outputName string
	var outputQuality int
	var threshold int

	flag.StringVar(&method, "m", "o4", "Dithering method (o4, o9, t, r)")
	flag.StringVar(&outputFormat, "f", "jpg", "Output file format (jpg, png)")
	flag.StringVar(&outputName, "o", "output", "Output file name")
	flag.IntVar(&outputQuality, "q", 100, "Output image quality (1-100)")
	flag.IntVar(&threshold, "t", 0, "Threshold value for threshold dithering")
	flag.Parse()

	file, _ := os.Open(flag.Arg(0))
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	draw.Draw(grayImg, bounds, img, bounds.Min, draw.Src)
	w, h := bounds.Max.X, bounds.Max.Y
	d := &Dither{
		sourceImage:   grayImg,
		width:         w,
		height:        h,
		newImage:      image.NewRGBA(bounds),
		outputFormat:  outputFormat,
		outputName:    outputName,
		outputQuality: outputQuality,
		threshold:     threshold,
	}

	switch method {
	case "o4":
		d.OrderedDither4()
	case "o9":
		d.OrderedDither9()
	case "t":
		d.ThresholdDither()
	case "r":
		d.RandomDither()
	case "a":
		d.outputName = fmt.Sprintf("%s_o4", outputName)
		d.OrderedDither4()
		d.SaveFile()
		d.outputName = fmt.Sprintf("%s_o9", outputName)
		d.OrderedDither9()
		d.SaveFile()
		d.outputName = fmt.Sprintf("%s_t", outputName)
		d.ThresholdDither()
		d.SaveFile()
		d.outputName = fmt.Sprintf("%s_r", outputName)
		d.RandomDither()
		d.SaveFile()

	default:
		log.Fatal("Unknown mode flag")
	}
}
