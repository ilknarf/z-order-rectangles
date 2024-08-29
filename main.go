package main

/**
 * Z-Order rectangle matching
 *
 * This code demonstrates usage of a Z-order representation of a rectangle to both store and query for
 * specific two-dimensional shapes.
 *
 * Z-order basically works by interweaving bits of a value into a larger value container.
 * For example, in order to weave two dimensions, the Z-order value of of (x, y) = (0b001, 0b110)
 * would be 0b010110 if the bits of x are placed first. If y's bits are placed first, then the
 * resulting value would be 0b101001.
 *
 * Why is this useful? This value effectively represents a multidimensional vector as a one-dimensional value,
 * In a form that is friendly to conventional sorting and indexing operations. Because the most significant
 * bits of all dimensions are placed towards the front, this means that n-dimensional structures can be
 * *approximately* queried efficiently without needing to use more "niche" data structures like quad trees.
 *
 * This particular script generates an image that demonstrates the bounds of using X bits of significance
 * when querying for a similar shape using a Z-Order representation of the shape.
 *
 * The output image contains three overlaid rectangles: Red = max spanning rectangle, Blue = min spanning
 * rectangle. The box outline by the top left corners and bottom right corners of the two rectangles designate
 * rectangles that will be matched using X bits of significance. The green rectangle is the original input.
 *
 * Note: I avoided negative values to avoid having to deal with two's complement edge cases.
 * Use cases that require this may need to add careful logic to account for such.
 **/

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"strconv"
)

const (
	x0 = 20.1234123
	y0 = 15.122122
	x1 = 89.99
	y1 = 57.999988

	width = 720

	height = 1080

	// number of bits within each dimension to use for precision
	// the less bits there are, the bigger the matches are.
	precisionBits = 4

	inc = 1 << (maxDimBits - precisionBits)

	fileName = "rect.jpg"
)

var (
	// some math to help calculate image
	gridUnitRatio = float32(inc) / maxDimVal
	hInc          = int(gridUnitRatio * height)
	wInc          = int(gridUnitRatio * width)
)

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// set background to white
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			img.Set(x, y, color.White)
		}
	}

	rect := &RectHash{}

	rect.SetX0(Int(x0))
	rect.SetX1(Int(x1))

	rect.SetY0(Int(y0))
	rect.SetY1(Int(y1))

	// calculate max and min range using rect val for convenience
	// truncating 4 bits for every dimension bit we want to truncate
	truncBits := (maxDimBits - precisionBits) * 4
	minRectVal := rect.Val >> uint64(truncBits)
	minRectVal <<= uint64(truncBits)

	minRect := &RectHash{Val: minRectVal}

	fmt.Printf("encoded value: %d\n", rect.Val)
	fmt.Printf("encoded binary value: 0b%s\n", strconv.FormatUint(rect.Val, 2))
	fmt.Println(rect.X0(), rect.X1(), rect.Y0(), rect.Y1())
	fmt.Printf("coords: x0 %f x1 %f y0 %f y1 %f\n", Ratio(rect.X0()), Ratio(rect.X1()), Ratio(rect.Y0()), Ratio(rect.Y1()))
	fmt.Printf("increment grid by w: %d h: %d\n", wInc, hInc)

	// calculate min and max span rectangles using above values
	minSpanRect := &RectHash{}
	minSpanRect.SetX0(minRect.X0() + inc)
	minSpanRect.SetY0(minRect.Y0() + inc)
	minSpanRect.SetX1(minRect.X1())
	minSpanRect.SetY1(minRect.Y1())

	maxSpanRect := &RectHash{}
	maxSpanRect.SetX0(minRect.X0())
	maxSpanRect.SetY0(minRect.Y0())
	maxSpanRect.SetX1(minRect.X1() + inc)
	maxSpanRect.SetY1(minRect.Y1() + inc)

	fmt.Println("drawing rectangles")
	// draw captured frame range and grid
	DrawRect(img, maxSpanRect, color.RGBA{
		R: 255,
		A: 0,
	})

	DrawRect(img, rect, color.RGBA{
		G: 255,
	})

	DrawRect(img, minSpanRect, color.RGBA{
		B: 255,
		A: 0,
	})

	fmt.Println("drawing grid")
	DrawGrid(img)
	SaveToFile(img, fileName)
}

func DrawRect(img *image.RGBA, rect *RectHash, col color.Color) {
	xstart := ToPx(rect.X0(), wInc)
	xend := ToPx(rect.X1(), wInc)

	ystart := ToPx(rect.Y0(), hInc)
	yend := ToPx(rect.Y1(), hInc)

	fmt.Println(xstart, xend, ystart, yend)

	for x := xstart; x < xend; x += 1 {
		for y := ystart; y < yend; y += 1 {
			img.Set(x, y, col)
		}
	}
}

func ToPx(vh uint, mult int) int {
	return int(float32(vh) / inc * float32(mult))
}

func DrawGrid(img *image.RGBA) {
	for x := 0; x < width; x += wInc {
		for y := 0; y < height; y += 1 {
			img.Set(x, y, color.Black)
		}
	}

	for y := 0; y < height; y += hInc {
		for x := 0; x < width; x += 1 {
			img.Set(x, y, color.Black)
		}
	}

}

func SaveToFile(img image.Image, path string) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = jpeg.Encode(f, img, nil)
	if err != nil {
		panic(err)
	}
}

func Ratio(val uint) float32 {
	return float32(val) / maxDimVal
}

func Int(val float32) uint {
	return uint((val / 100) * maxDimVal)
}
