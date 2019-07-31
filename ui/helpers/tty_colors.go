package helpers

import (
	commanderWidgets "github.com/daylioti/docker-commander/ui/widgets"
	"github.com/gizak/termui/v3"
	"strconv"
	"unicode"
)

// StyleInit - rune for "["
const StyleInit = 91

// StyleEnd - rune for "m"
const StyleEnd = 109

// StyleAttrLength - means length of "[m"
const StyleAttrLength = 2

// StyleMinLength - means length of ["2"m
const StyleMinLength = 1

// DefaultColor - back to default color rune
var resetStyle = []byte{39, 0, 1, 2, 4, 5, 7, 8, 22}

// vt100Codes16 codes from https://misc.flogisoft.com/bash/tip_colors_and_formatting
var vt100Codes16 = map[byte]string{
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

// TTYColorsParse - convert ANSI colors to termui.
func TTYColorsParse(buff []byte) []byte {
	var index int
	var styleByte byte
	var styleLength int
	var step int
	var exist bool
	var style string
	var replace []byte
	for {
		if len(buff) < index+StyleMinLength+StyleAttrLength || len(buff) <= 6 || index < 0 {
			return buff
		}
		if len(buff) > index+StyleAttrLength && buff[index] == StyleInit && unicode.IsNumber(rune(buff[index+1])) {
			// Possible style.
			if len(buff) > index+StyleAttrLength {
				if buff[index+StyleAttrLength] == StyleEnd {
					styleByte = convertStyleByte([]byte{buff[index+1]})
					styleLength = 1
				}
				if unicode.IsNumber(rune(buff[index+2])) {
					if len(buff) > index+3 && buff[index+3] == StyleEnd {
						styleByte = convertStyleByte([]byte{buff[index+1], buff[index+2]})
						styleLength = 2
					} else if len(buff) > index+4 && unicode.IsNumber(rune(buff[index+3])) && buff[index+4] == StyleEnd {
						styleByte = convertStyleByte([]byte{buff[index+1], buff[index+2], buff[index+3]})
						styleLength = 3
					}
				}
				if !isReset(styleByte) {
					style, exist = vt100Codes16[styleByte]
					if exist && len(buff) > index+styleLength+7 {
						// 7 means minimum length of closing style.
						// Replace start ANSI style to termui start style byte.
						replace = []byte(string(commanderWidgets.TokenBeginStyledText))
						buff = append(buff[:index+styleLength-StyleAttrLength], append(replace, buff[index+styleLength+StyleAttrLength:]...)...)
						index += len(replace) - styleLength
					} else {
						// Style not in list, remove them.
						if len(buff) < index+styleLength+StyleAttrLength {
							step = 1
						} else {
							step = 2
						}
						buff = append(buff[:index], buff[index+styleLength+step:]...)
						index -= styleLength + step
					}
				} else {
					if exist {
						replace = []byte(string(commanderWidgets.TokenEndStyledText) + string(commanderWidgets.TokenBeginStyle) + style + string(commanderWidgets.TokenEndStyle))
						// Paste termui style.
						buff = append(buff[:index+styleLength-StyleAttrLength], append(replace, buff[index+styleLength+StyleAttrLength:]...)...)
						index += len(replace)
					} else {
						// Remove reset chars.
						if len(buff) < index+styleLength+StyleAttrLength {
							step = 1
						} else {
							step = 2
						}
						buff = append(buff[:index], buff[index+styleLength+step:]...)
						index -= styleLength + step
					}
					style = ""
					exist = false
				}
			}
		}
		index++
	}
}

// isReset check for closing style byte.
func isReset(checkByte byte) bool {
	for _, reset := range resetStyle {
		if checkByte == reset {
			return true
		}
	}
	return false
}

// convertStyleRune convert color code to int8 format.
func convertStyleByte(style []byte) byte {
	var result int64
	result, _ = strconv.ParseInt(string(style), 10, 8)
	return byte(result)
}
