package pixelizer

import (
	"image"
	"sync"
)

type Pixelizer interface {
	Pixelize(image.Image) image.Image
}

type simplePixelizer struct {
	blockSize  int
	colorRange int
}

func NewSimplePixelizer(blockSize int) Pixelizer {
	return &simplePixelizer{
		blockSize:  blockSize,
		colorRange: blockSize / 2,
	}
}

func (p *simplePixelizer) Pixelize(i image.Image) image.Image {
	bounds := i.Bounds()

	resX, resY := bounds.Max.X, bounds.Max.Y
	if resX%p.blockSize != 0 {
		resX -= resX % p.blockSize
	}
	if resY%p.blockSize != 0 {
		resY -= resY % p.blockSize
	}

	pi := image.NewNRGBA(image.Rect(0, 0, resX, resY))

	var wg sync.WaitGroup
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + p.blockSize {
		wg.Add(1)

		go func(y int) {
			defer wg.Done()

			for x := bounds.Min.X; x < bounds.Max.X; x = x + p.blockSize {
				c := i.At(x+p.colorRange, y+p.colorRange)
				for rx := x; rx < x+p.blockSize; rx++ {
					for ry := y; ry < y+p.blockSize; ry++ {
						pi.Set(rx, ry, c)
					}
				}
			}
		}(y)
	}
	wg.Wait()

	return pi
}
