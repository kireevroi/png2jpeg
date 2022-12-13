package main

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func getNames(path string) []string {
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

func png2jpeg(path string, name string) {
	pngImgFile, err := os.Open(path + name)
	if err != nil {
		log.Fatal(err)
	}
	defer pngImgFile.Close()
	imgSrc, err := png.Decode(pngImgFile)
	if err != nil {
		log.Fatal(err)
	}
	newImg := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)
	jpgImgFile, err := os.Create("./jpeg/" + strings.TrimSuffix(name, filepath.Ext(name)) + ".jpeg")
	if err != nil {
		log.Fatal(err)
	}
	defer jpgImgFile.Close()
	var opt jpeg.Options
	opt.Quality = 80
	err = jpeg.Encode(jpgImgFile, newImg, &opt)
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
}

func squeezeSingleJpeg(path string, name string) {
	file, err := os.Open(path + name)
	if err != nil {
		log.Fatal(err)
	}
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	m := resize.Resize(256, 0, img, resize.Lanczos3)
	out, err := os.Create("squeeze/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	jpeg.Encode(out, m, nil)
	wg.Done()
}

func СonvertBatchPng(path string, arr []string) {
	if _, err := os.Stat("jpeg"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("jpeg", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	for _, name := range arr {
		wg.Add(1)
		go png2jpeg(path, name)
	}
	wg.Wait()
}

func SqueezeBatchJpeg(path string, arr []string) {
	if _, err := os.Stat("squeeze"); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir("squeeze", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
	for _, name := range arr {
		wg.Add(1)
		go squeezeSingleJpeg(path, name)
	}
	wg.Wait()
}

func main() {
	СonvertBatchPng("pokemon/", getNames("pokemon/"))
	SqueezeBatchJpeg("jpeg/", getNames("jpeg/"))
}
