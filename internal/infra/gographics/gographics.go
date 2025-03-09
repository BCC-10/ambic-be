package gographics

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type GoGraphicsIf interface {
	GenerateBratImage(width, height int, text string, bgColorOpt, textColorOpt *string) (*multipart.FileHeader, error)
}

type GoGraphics struct {
}

func NewGoGraphics() GoGraphicsIf {
	return &GoGraphics{}
}

const fontURL = "https://uqrnpvnydqwyreivieki.supabase.co/storage/v1/object/public/uploads/fonts/Poppins-Medium.ttf"

func downloadFont(url string) (*truetype.Font, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download font: %v", err)
	}
	defer resp.Body.Close()

	fontData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read font data: %v", err)
	}

	font, err := truetype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
	}

	return font, nil
}

func parseColor(hex string, defaultColor color.Color) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return defaultColor
	}
	if r, err := strconv.ParseInt(hex[0:2], 16, 32); err == nil {
		if g, err := strconv.ParseInt(hex[2:4], 16, 32); err == nil {
			if b, err := strconv.ParseInt(hex[4:6], 16, 32); err == nil {
				return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
			}
		}
	}
	return defaultColor
}

func (g *GoGraphics) GenerateBratImage(width, height int, text string, bgColorOpt, textColorOpt *string) (*multipart.FileHeader, error) {
	dc := gg.NewContext(width, height)

	// Parse background color (default: RGB(62, 151, 149))
	bgColor := "#3E9795"
	if bgColorOpt != nil {
		bgColor = *bgColorOpt
	}
	bg := parseColor(bgColor, color.RGBA{R: 62, G: 151, B: 149, A: 255})
	dc.SetColor(bg)
	dc.Clear()

	// Download font
	font, err := downloadFont(fontURL)
	if err != nil {
		return nil, err
	}

	// Calculate font size based on image size
	fontSize := float64(height) * 0.2
	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	dc.SetFontFace(face)

	// Parse text color (default: white)
	textColor := "#FFFFFF"
	if textColorOpt != nil {
		textColor = *textColorOpt
	}
	txtColor := parseColor(textColor, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.SetColor(txtColor)

	// Draw the text centered
	dc.DrawStringAnchored(text, float64(width)/2, float64(height)/2, 0.5, 0.5)

	// Encode to buffer
	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	// Convert buffer to multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: "brat.png",
		Size:     int64(buf.Len()),
	}

	return fileHeader, nil
}
