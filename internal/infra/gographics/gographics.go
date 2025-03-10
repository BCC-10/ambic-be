package gographics

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

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

func generateBratImage(width, height int, text string, bgColorOpt, textColorOpt *string) (*multipart.FileHeader, io.Reader, error) {
	dc := gg.NewContext(width, height)

	bgColor := "#3E9795"
	if bgColorOpt != nil {
		bgColor = *bgColorOpt
	}
	bg := parseColor(bgColor, color.RGBA{R: 62, G: 151, B: 149, A: 255})
	dc.SetColor(bg)
	dc.Clear()

	font, err := downloadFont(fontURL)
	if err != nil {
		return nil, nil, err
	}

	textColor := "#FFFFFF"
	if textColorOpt != nil {
		textColor = *textColorOpt
	}
	txtColor := parseColor(textColor, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	dc.SetColor(txtColor)

	maxWidth := float64(width) * 0.9
	maxHeight := float64(height) * 0.9
	fontSize := float64(height) * 0.2

	var lines []string
	for fontSize > 10 {
		face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
		dc.SetFontFace(face)
		lines = dc.WordWrap(text, maxWidth)

		totalTextHeight := float64(len(lines)) * (fontSize * 1.2)

		if totalTextHeight <= maxHeight {
			break
		}
		fontSize -= 2
	}

	totalTextHeight := float64(len(lines)) * (fontSize * 1.2)

	startY := (float64(height)-totalTextHeight)/2 + fontSize/3

	for i, line := range lines {
		y := startY + float64(i)*(fontSize*1.2)
		dc.DrawStringAnchored(line, float64(width)/2, y, 0.5, 0.5)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, nil, fmt.Errorf("failed to encode image: %v", err)
	}

	fileHeader := &multipart.FileHeader{
		Filename: "brat.png",
		Size:     int64(buf.Len()),
	}

	return fileHeader, bytes.NewReader(buf.Bytes()), nil
}

func NewImage(width, height int, text string) (*multipart.FileHeader, io.Reader) {
	fileHeader, fileReader, err := generateBratImage(width, height, text, nil, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil
	}

	return fileHeader, fileReader
}
