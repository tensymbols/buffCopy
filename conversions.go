package main

import (
	"bufio"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type UpperCase struct {
	output io.Writer
}
type LowerCase struct {
	output io.Writer
}
type TrimSpaces struct {
	input    *bufio.Reader
	output   io.Writer
	limit    uint
	SpaceBuf []rune
}

func NewUpperCase(output io.Writer) io.Writer {
	return &UpperCase{output}
}
func NewLowerCase(output io.Writer) io.Writer {
	return &LowerCase{output}
}
func NewTrimSpaces(input *bufio.Reader, output io.Writer, limit uint) TrimSpaces {
	return TrimSpaces{input, output, limit, []rune{}}
}

func (uc *UpperCase) Write(p []byte) (int, error) {
	data := []byte(strings.ToUpper(string(p)))

	return uc.output.Write(data)
}
func (lc *LowerCase) Write(p []byte) (int, error) {
	data := []byte(strings.ToLower(string(p)))
	return lc.output.Write(data)
}

func (ts *TrimSpaces) Write(p []byte) (int, error) {
	strBuf := string(p)
	p = []byte{}
	for _, v := range strBuf {
		if unicode.IsSpace(v) {
			ts.SpaceBuf = append(ts.SpaceBuf, v)
		} else {
			for _, vv := range ts.SpaceBuf {
				p = utf8.AppendRune(p, vv)
			}
			ts.SpaceBuf = []rune{}

			p = utf8.AppendRune(p, v)
		}
	}
	return ts.output.Write(p)
}
func (ts *TrimSpaces) Read(p []byte) (int, error) {
	var cnt int
	r, n, err := ts.input.ReadRune()
	cnt += n
	if err != nil {

		return 0, err
	}
	for unicode.IsSpace(r) && cnt < int(ts.limit) && err == nil {
		r, n, err = ts.input.ReadRune()
		cnt += n
	}

	if !unicode.IsSpace(r) {
		p = []byte(string(p) + string(r))
		_, err = ts.Write(p)

	}
	return cnt, err
}
