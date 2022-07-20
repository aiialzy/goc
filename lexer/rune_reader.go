package lexer

import (
	"errors"
	"io"
	"unicode/utf8"
)

type runeReader struct {
	runes []rune
	index int
}

func (rr *runeReader) ReadRune() (rune, int, error) {
	if rr.index >= len(rr.runes) {
		return 0, 0, io.EOF
	}
	r := rr.runes[rr.index]
	rr.index++
	return r, utf8.RuneLen(r), nil
}

func (rr *runeReader) UnreadRune() error {
	if rr.index == 0 {
		return errors.New("min index is 0")
	}

	rr.index--
	return nil
}

func (rr *runeReader) PeekRune(t int) (rune, error) {
	if t < 0 {
		return 0, errors.New("must be bigger than or equal 0")
	}

	if rr.index+t >= len(rr.runes) {
		return 0, io.EOF
	}

	return rr.runes[rr.index+t], nil
}
