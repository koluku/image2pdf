package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"

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

	if flag.NArg() == 0 {
		log.Fatalf("%+v", errors.WithStack(errors.New("input directory is empty")))
	}

	for _, inputDir := range flag.Args() {
		paths, err := os.ReadDir(inputDir)
		if err != nil {
			log.Fatalf("%+v", errors.WithStack(err))
		}
		absPath, err := filepath.Abs(inputDir)
		if err != nil {
			log.Fatalf("%+v", errors.WithStack(err))
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

			file, err := os.Open(filepath.Join(absPath, path.Name()))
			if err != nil {
				log.Fatalf("%+v", errors.WithStack(err))
			}
			defer file.Close()

			config, _, err := image.DecodeConfig(file)
			if err != nil {
				log.Fatalf("%+v", errors.WithStack(err))
			}

			if config.Height >= config.Width {
				pdf.AddPageWithOption(gopdf.PageOption{
					PageSize: &gopdf.Rect{
						W: float64(842) * float64(config.Width) / float64(config.Height),
						H: 842,
					},
				})
				pdf.Image(filepath.Join(absPath, path.Name()), 0, 0, &gopdf.Rect{
					W: float64(842) * float64(config.Width) / float64(config.Height),
					H: 842,
				})
			} else {
				pdf.AddPageWithOption(gopdf.PageOption{
					PageSize: &gopdf.Rect{
						W: 842,
						H: float64(842) * float64(config.Height) / float64(config.Width),
					},
				})
				pdf.Image(filepath.Join(absPath, path.Name()), 0, 0, &gopdf.Rect{
					W: 842,
					H: float64(842) * float64(config.Height) / float64(config.Width),
				})
			}
		}
		splitedDirs := strings.Split(filepath.ToSlash(absPath), "/")
		if err := pdf.WritePdf(filepath.Join(absPath, fmt.Sprintf("%s.pdf", splitedDirs[len(splitedDirs)-1]))); err != nil {
			log.Fatalf("%+v", errors.WithStack(err))
		}
	}
}
