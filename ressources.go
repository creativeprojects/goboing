package main

import (
	"image"

	_ "image/png"

	"github.com/hajimehoshi/ebiten"
	"github.com/markbates/pkger"
)

func loadImages(imageNames []string) (map[string]*ebiten.Image, error) {
	imagesMap := make(map[string]*ebiten.Image, len(imageNames))
	for _, imageName := range imageNames {
		file, err := pkger.Open("/images/" + imageName)
		if err != nil {
			return imagesMap, err
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			return imagesMap, err
		}
		img2, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
		if err != nil {
			return imagesMap, err
		}
		imagesMap[imageName] = img2
	}
	return imagesMap, nil
}
