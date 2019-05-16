package helpers

import (
	"github.com/gizak/termui/v3"
	"strconv"
	"unicode"
)

const StyleInit = 91 // [
const StyleEnd = 109 // m
const DefaultColor = 39

// vt100Codes16 codes from https://misc.flogisoft.com/bash/tip_colors_and_formatting
var vt100Codes16 = map[rune]string{
	30: "fg:black",
	31: "fg:red",
	32: "fg:green",
	33: "fg:yellow",
	34: "fg:blue",
	35: "fg:magenta",
	36: "fg:cyan",
	37: "fg:254",
	90: "fg:244",
	91: "fg:160",
	92: "fg:46",
	93: "fg:226",
	94: "fg:111",
	95: "fg:134",
	96: "fg:87",
	97: "fg:white",

	49: "bg:white",
	40: "bg:black",
	41: "bg:red",
	42: "bg:green",
	43: "bg:yellow",
	44: "bg:blue",
	45: "bg:magenta",
	46: "bg:cyan",
	47: "bg:254",
}

// GetAllTermColors return all colors that termui can render.
func GetAllTermColors() map[string]termui.Color {
	colors := termui.StyleParserColorMap
	for i := 0; i < 255; i++ {
		colors[strconv.Itoa(i)] = termui.Color(i)
	}
	return colors
}

// TTYColorsParse replace tty colors to termui.
func TTYColorsParse(buffer []rune) []rune {
	var style string
	var styleRune rune
	var exist bool
	for bufferIndex := 0; bufferIndex < len(buffer); bufferIndex++ {
		if buffer[bufferIndex] == StyleInit {
			switch {
			case !exist && len(buffer) > bufferIndex+1 && buffer[bufferIndex+1] == StyleInit:
				// Remove [
				buffer = append(buffer[:bufferIndex], buffer[bufferIndex+1:]...)

			case len(buffer) > bufferIndex+3 && buffer[bufferIndex+3] == StyleEnd:
				styleRune = convertStyleRune([]rune{buffer[bufferIndex+1], buffer[bufferIndex+2]})
				if styleRune == DefaultColor {
					if exist {
						buffer = append(buffer[:bufferIndex], append([]rune("]("+style+")"), buffer[bufferIndex+4:]...)...)
					} else {
						buffer = append(buffer[:bufferIndex], buffer[bufferIndex+4:]...)
					}
				} else {
					buffer = append(buffer[:bufferIndex+1], buffer[bufferIndex+4:]...)
				}
				style, exist = vt100Codes16[styleRune]
				if !exist && styleRune != DefaultColor && len(buffer) > bufferIndex+1 {
					buffer = append(buffer[:bufferIndex], buffer[bufferIndex+1:]...)
				}

			case len(buffer) > bufferIndex+3 && buffer[bufferIndex+2] == StyleEnd && unicode.IsNumber(buffer[bufferIndex+1]):
				buffer = append(buffer[:bufferIndex], buffer[bufferIndex+3:]...)

			default:
				buffer = append(buffer[:bufferIndex], buffer[bufferIndex+1:]...)
			}
		} else if exist && buffer[bufferIndex] == ']' {
			buffer = append(buffer[:bufferIndex], buffer[bufferIndex+1:]...)
		}
	}
	return buffer
}

// convertStyleRune convert color code to int32 format.
func convertStyleRune(style []rune) rune {
	var result int64
	result, _ = strconv.ParseInt(string(style), 10, 32)
	return int32(result)
}
