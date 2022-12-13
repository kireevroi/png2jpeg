package main

import (
	"fmt"
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

         // convert newImage to JPEG encoded byte and save to jpgImgFile
         // with quality = 80
    err = jpeg.Encode(jpgImgFile, newImg, &opt)

         //err = jpeg.Encode(jpgImgFile, newImg, nil) -- use nil if ignore quality options

    if err != nil {
      log.Println(err)
      os.Exit(1)
    }
    fmt.Println("Converted PNG file to JPEG file")
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



         // create image from PNG file


         // create a new Image with the same dimension of PNG image
 }

 func main() {
	convertPng2Jpeg("pokemon/", getNames("pokemon/"))
 }
