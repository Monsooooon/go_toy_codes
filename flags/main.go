package main

import (
	"flag"
	"fmt"
	"image/color"
	"strconv"
)

type colorValue struct {
	color.Color // struct embedding
}

func (c *colorValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return fmt.Errorf("input is not a color: %v", err)
	}
	b := uint8(v & 0xFF)
	g := uint8((v >> 8) & 0xFF)
	r := uint8((v >> 16) & 0xFF)

	c.Color = color.RGBA{R: r, G: g, B: b, A: 0xFF}
	return nil
}

func (c *colorValue) String() string {
	var r, g, b, a uint32
	if c != nil && c.Color != nil {
		r, g, b, a := c.RGBA()
		r, g, b, a = r/256, g/256, b/256, a/256
	}
	return fmt.Sprintf("rgba(%v %v %v %v)", r, g, b, a)
}

func main() {

	var fg, bg colorValue
	flag.Var(&fg, "fg", "foreground color")
	flag.Var(&bg, "bg", "background color")

	flag.Parse()
	draw(&fg, &bg)
}

func draw(fg, bg color.Color) {
	fmt.Printf("drawing with foreground color %v and background color %v", fg, bg)
}
