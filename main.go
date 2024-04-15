package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/signintech/gopdf"
)

var (
	author   string
	dest     string
	noRotate bool
)

func init() {
	flag.StringVar(&author, "author", "", "author")
	flag.StringVar(&dest, "dest", "", "dest")
	flag.BoolVar(&noRotate, "no-rotate", false, "noRotate")
}

func main() {
	flag.Parse()

	if author == "" {
		log.Fatal("--author is empty")
	}

	if dest == "" {
		log.Fatal("--dest is empty")
	}

	if flag.NArg() < 1 {
		log.Fatal("input directory is empty")
	}

	if err := run(flag.Args()); err != nil {
		log.Fatal(err)
	}
}

func run(inputDirs []string) error {
	for _, inputDir := range inputDirs {
		doc, err := NewDoc(inputDir)
		if err != nil {
			return err
		}
		if err := doc.toPDF(); err != nil {
			return err
		}
	}

	return nil
}

type Doc struct {
	Author     string
	Title      string
	ImageFiles []string
}

func NewDoc(inputDir string) (*Doc, error) {
	title := filepath.Base(inputDir)

	entries, err := os.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	imageFiles := newImageFiles(inputDir, entries)

	return &Doc{
		Author:     author,
		Title:      title,
		ImageFiles: imageFiles,
	}, nil
}

const A4_HEIGHT = 4093

func (d *Doc) toPDF() error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		Unit: gopdf.Unit_PT,
	})

	for _, imageFilePath := range d.ImageFiles {
		file, err := os.Open(imageFilePath)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		defer file.Close()

		config, _, err := image.DecodeConfig(file)
		if err != nil {
			log.Fatalf("%+v", err)
		}

		var width, height float64
		if config.Width > config.Height && noRotate {
			width = A4_HEIGHT
			height = float64(A4_HEIGHT) * float64(config.Height) / float64(config.Width)
		} else {
			width = float64(A4_HEIGHT) * float64(config.Width) / float64(config.Height)
			height = A4_HEIGHT
		}
		pdf.AddPageWithOption(gopdf.PageOption{
			PageSize: &gopdf.Rect{
				W: width,
				H: height,
			},
		})
		pdf.Image(imageFilePath, 0, 0, &gopdf.Rect{
			W: width,
			H: height,
		})
	}

	filename := fmt.Sprintf("%s_%s.pdf", d.Author, d.Title)
	if err := pdf.WritePdf(filepath.Join(dest, filename)); err != nil {
		return err
	}

	return nil
}

func newImageFiles(parentDir string, entries []fs.DirEntry) []string {
	imagePathes := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !isImageFile(entry.Name()) {
			continue
		}
		imagePathes = append(imagePathes, filepath.Join(parentDir, entry.Name()))
	}
	return imagePathes
}

var exts = []string{".jpeg", ".jpg", ".png"}

func isImageFile(name string) bool {
	for _, ext := range exts {
		if filepath.Ext(name) == ext {
			return true
		}
	}
	return false
}
