package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/czertbytes/pixelizer2/pixelizer"
)

func main() {
	outFile := flag.String("output", "pix.png", "The output pixelized file.")
	blockSize := flag.Int("block-size", 8, "The pixel block size.")
	flag.Parse()
	flag.Usage = func() {
		fmt.Printf("\nUsage: pixelize [options] <file>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file := ""
	for _, p := range flag.Args() {
		if strings.HasPrefix(p, "-") {
			continue
		}
		file = p
	}

	start := time.Now()

	in, err := os.Open(file)
	if err != nil {
		fmt.Printf("Reading input file %q failed, error: %s", in, err)
		os.Exit(1)
	}
	defer in.Close()

	inImage, _, err := image.Decode(in)
	if err != nil {
		fmt.Printf("Reading input file %q failed, error: %s", in, err)
		os.Exit(1)
	}

	outImage := pixelizer.NewSimplePixelizer(*blockSize).Pixelize(inImage)

	out, err := os.Create(*outFile)
	if err != nil {
		fmt.Printf("Writing output image %q failed, error: %s", *outFile, err)
		os.Exit(1)
	}
	defer out.Close()

	switch strings.ToLower(filepath.Ext(out.Name())) {
	case ".jpg", ".jpeg":
		if err := jpeg.Encode(out, outImage, &jpeg.Options{Quality: 80}); err != nil {
			fmt.Printf("Encoding output image %q failed, error: %s", *outFile, err)
			os.Exit(1)
		}
	case ".png":
		fallthrough
	default:
		if err := png.Encode(out, outImage); err != nil {
			fmt.Printf("Encoding output image %q failed, error: %s", *outFile, err)
			os.Exit(1)
		}
	}

	fmt.Printf("Done in %s\n", time.Now().Sub(start))

}
