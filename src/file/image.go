package file

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
)

// PngToJpeg png to jpeg
// quality: -1 noting
func PngToJpeg(img image.Image, jpegF io.Writer, quality int) error {
	// create a new Image with the same dimension of PNG image
	newImg := image.NewRGBA(img.Bounds())

	// we will use white background to replace PNG's transparent background
	// you can change it to whichever color you want with
	// a new color.RGBA{} and use image.NewUniform(color.RGBA{<fill in color>}) function
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// paste PNG image OVER to newImage
	draw.Draw(newImg, newImg.Bounds(), img, img.Bounds().Min, draw.Over)

	var q *jpeg.Options
	if quality > 0 {
		q = &jpeg.Options{Quality: quality}
	}
	err := jpeg.Encode(jpegF, newImg, q)
	if err != nil {
		return err
	}

	return nil
}

// ContainsTransparent judge png contains transparent
func ContainsTransparent(img image.Image) bool {
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			rgb := img.At(i, j)
			_, _, _, a := rgb.RGBA()
			if a == 0 { // transparent
				return true
			}
		}
	}
	return false
}

// ImageType jpeg,png,gif
func ImageType(buffer []byte) string {
	contentType := ContentType(buffer)
	switch contentType {
	case "image/jpeg":
		return JPEG
	case "image/png":
		return PNG
	case "image/gif":
		return GIF
	default:
		return ""
	}
}
