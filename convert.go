package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"errors"
	"github.com/nfnt/resize"
)

 func getNames(path string) []string{
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	list, err := file.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}
	arr := make([]string, 0)

	for _, f := range list {
		arr = append(arr, f.Name())
	}
	return arr
 }

func convertSinglePng2Jpeg(path string, name string) {
		pngImgFile, err := os.Open(path + name)
    if err != nil {
      log.Println(err)
      os.Exit(1)
    }
    defer pngImgFile.Close()
    imgSrc, err := png.Decode(pngImgFile)
    if err != nil {
      log.Println(err)
      os.Exit(1)
    }
    newImg := image.NewRGBA(imgSrc.Bounds())
    draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
    draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)
    jpgImgFile, err := os.Create("./jpeg/" + strings.TrimSuffix(name, filepath.Ext(name)) + ".jpeg")

    if err != nil {
      log.Println(err)
      os.Exit(1)
    }
    defer jpgImgFile.Close()
    var opt jpeg.Options
    opt.Quality = 80
    err = jpeg.Encode(jpgImgFile, newImg, &opt)
    if err != nil {
      log.Println(err)
      os.Exit(1)
    }
    log.Println("Converted PNG file to JPEG file")
}

func SqueezeBatchJpeg(path string, arr []string) {
	if _, err := os.Stat("squeeze"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("squeeze", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	for _, name := range arr {
    SqueezeSingleJpeg(path, name)
	}
}

func SqueezeSingleJpeg(path string, name string) {
	// open "test.jpg"
	file, err := os.Open(path + name)
	if err != nil {
		log.Fatal(err)
	}
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(256, 0, img, resize.Lanczos3)

	out, err := os.Create("squeeze/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)
	log.Println("Squeezed JPEG")
}

 func convertPng2Jpeg(path string, arr []string) {
	if _, err := os.Stat("jpeg"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("jpeg", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	for _, name := range arr {
    convertSinglePng2Jpeg(path, name)
	}
 }

 func main() {
	convertPng2Jpeg("pokemon/", getNames("pokemon/"))
	SqueezeBatchJpeg("jpeg/", getNames("jpeg/"))
 }
