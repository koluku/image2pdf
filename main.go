package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/signintech/gopdf"
)

var exts = []string{".jpeg", ".jpg", ".png"}

func isImageFile(name string) bool {
	for _, ext := range exts {
		if filepath.Ext(name) == ext {
			return true
		}
	}
	return false
}

func main() {
	flag.Parse()

	inputDir := flag.Arg(0)
	if inputDir == "" {
		log.Fatal("No argument error")
	}

	paths, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		Unit: gopdf.Unit_PT,
	})

	for _, path := range paths {
		if path.IsDir() {
			continue
		}
		if !isImageFile(path.Name()) {
			continue
		}

		file, _ := os.Open(path.Name())
		defer file.Close()

		src, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}

		pdf.AddPageWithOption(gopdf.PageOption{
			PageSize: &gopdf.Rect{
				W: float64(src.Bounds().Dx()),
				H: float64(src.Bounds().Dy()),
			},
		})
		pdf.Image(path.Name(), 0, 0, &gopdf.Rect{
			W: float64(src.Bounds().Dx()),
			H: float64(src.Bounds().Dy()),
		})
	}

	ab, err := filepath.Abs(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	splitedDirs := strings.Split(filepath.ToSlash(ab), "/")
	if err := pdf.WritePdf(fmt.Sprintf("%s.pdf", splitedDirs[len(splitedDirs)-1])); err != nil {
		log.Fatal(err)
	}
}
