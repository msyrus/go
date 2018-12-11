package sms

import (
	"math"
	"strings"
	"unicode/utf8"
)

// This is the Golang port of https://github.com/messente/sms-length-calculator

type GSMCharSet int

const (
	GSMCharSetGSM7 GSMCharSet = iota
	GSMCharSetUnicode
)

const escRune = '\x1b'

func Len(str string) int {
	cset := Encoding(str)
	clen := utf8.RuneCountInString(str)
	if cset == GSMCharSetUnicode {
		if clen <= 70 {
			return 1
		}
		return int(math.Ceil(float64(clen) / 67))
	}

	for _, c := range gsm7bitCharsExt {
		str = strings.Replace(str, string(c), string([]rune{escRune, c}), -1)
	}

	clen = utf8.RuneCountInString(str)

	if clen <= 160 {
		return 1
	}

	scout := 0
	sruns := []rune(str)
	for len(sruns) > 0 {
		scout++
		if len(sruns) < 153 {
			break
		}
		if sruns[152] == escRune {
			sruns = sruns[152:]
			continue
		}
		sruns = sruns[153:]
	}

	return scout
}

func Encoding(str string) GSMCharSet {
	gsm7AllChars := append(gsm7BitChars, gsm7bitCharsExt...)
	cmap := map[rune]bool{}
	for _, c := range gsm7AllChars {
		cmap[c] = true
	}

	f := func(c rune) bool {
		return !cmap[c]
	}
	if strings.IndexFunc(str, f) == -1 {
		return GSMCharSetGSM7
	}
	return GSMCharSetUnicode
}

var gsm7BitChars = []rune{
	'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
	'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', '\x1b', 'Æ', 'æ', 'ß', 'É',
	' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
	'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
	'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
	'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
}

var gsm7bitCharsExt = []rune{
	'\f', '^', '{', '}', '\\', '[', '~', ']', '|', '€',
}
