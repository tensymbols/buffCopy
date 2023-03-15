package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"unicode/utf8"
)

type Options struct {
	From      string
	To        string
	Offset    uint
	Limit     uint
	BlockSize uint
	Conv      Conversions
}
type Conversions struct {
	UpperCase  bool
	LowerCase  bool
	TrimSpaces bool
}

func ValidOptions(opts *Options, params []string) error {
	if len(params) > 2 {
		return errors.New("invalid params")
	}
	for _, k := range params {
		switch k {
		case "upper_case":
			opts.Conv.UpperCase = true
		case "lower_case":
			opts.Conv.LowerCase = true
		case "trim_spaces":
			opts.Conv.TrimSpaces = true
		default:
			return errors.New("invalid params")
		}
	}
	if opts.Conv.LowerCase && opts.Conv.UpperCase {
		return errors.New("uppercase and lowercase simultaneously")
	}
	return nil
}

func Copy(rw *bufio.ReadWriter, opts *Options) error {
	var err error
	var buf bytes.Buffer
	var limitCnt uint = 0

	var caseConverter io.Writer
	var spaceTrimmer io.ReadWriter

	if opts.Conv.LowerCase {
		caseConverter = NewLowerCase(&buf)
	} else if opts.Conv.UpperCase {
		caseConverter = NewUpperCase(&buf)
	}
	if opts.Conv.TrimSpaces {
		var s = NewTrimSpaces(rw.Reader, &buf, opts.Limit)
		spaceTrimmer = &s
	}

	for i := uint(0); i < opts.Offset && err == nil; i++ {
		_, err = rw.ReadByte()
	}
	if err == io.EOF {
		return errors.New("offset>=input size")
	}

	if spaceTrimmer != nil {
		var n int
		n, err = spaceTrimmer.Read(buf.Bytes())
		limitCnt += uint(n)
	}

	for err == nil && limitCnt < opts.Limit {
		for i := uint(0); i < opts.BlockSize && err == nil && limitCnt < opts.Limit; i++ {

			var b byte
			b, err = rw.ReadByte()
			buf.WriteByte(b)

			limitCnt++
		}
		bufValid := false
		ixValid := 0

		for i := 0; i <= 3; i++ {
			bufValid = bufValid || utf8.Valid(buf.Bytes()[i:])
		}

		for !bufValid && limitCnt < opts.Limit && err == nil { // if current buffer isn't a valid utf8 string then append bytes until it becomes one
			var b byte
			b, err = rw.ReadByte()

			buf.WriteByte(b)
			bufValid = bufValid || utf8.Valid(buf.Bytes()[ixValid:])
			limitCnt++
		}

		if buf.Bytes()[buf.Len()-1] == '\x00' {
			buf.Truncate(buf.Len() - 1)
		}

		var writeErr error
		if caseConverter != nil {
			bufBytes := buf.Bytes()
			buf.Reset()
			_, writeErr = caseConverter.Write(bufBytes)
		}
		if spaceTrimmer != nil {
			bufBytes := buf.Bytes()
			buf.Reset()
			_, writeErr = spaceTrimmer.Write(bufBytes)
		}

		if writeErr != nil {
			return writeErr
		}

		_, err := rw.Write(buf.Bytes())
		if err != nil {
			return err
		}
		buf.Reset()

	}
	errFLush := rw.Writer.Flush()
	if errFLush != nil {
		return errFLush
	}
	return nil
}
func ParseFlags() (*Options, error) {
	var opts Options
	flag.StringVar(&opts.From, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.To, "to", "", "file to write. by default - stdout")
	flag.UintVar(&opts.Offset, "offset", 0, "offset in bytes. by default - 0")
	flag.UintVar(&opts.Limit, "limit", math.MaxInt, "limit in bytes. by default - from offset to the end)")
	flag.UintVar(&opts.BlockSize, "block-size", 1024, "block size in bytes. by default - 1024")

	convOptions := flag.String("conv", "", "limit in bytes. by default - no options")
	flag.Parse()
	var params []string
	if *convOptions != "" {
		opts.Conv = Conversions{false, false, false}
		params = strings.Split(*convOptions, ",")
		err := ValidOptions(&opts, params)
		if err != nil {
			return &opts, err
		}
	}
	if opts.BlockSize < 4 {
		opts.BlockSize = 4
	}
	return &opts, nil
}

func main() {
	opts, err := ParseFlags()

	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can't parse flags:", err)
		os.Exit(1)
	}

	var br *bufio.Reader
	var bw *bufio.Writer

	if opts.From == "" {
		br = bufio.NewReader(os.Stdin)
	} else {
		f, errOpen := os.Open(opts.From)
		if errOpen != nil {
			panic("can't open source file")
		}
		br = bufio.NewReader(f)

	}
	if opts.To == "" {
		bw = bufio.NewWriter(os.Stdout)
	} else {

		if _, errExists := os.Stat(opts.To); errExists == nil {
			panic("destination file already exists")
		}
		f, errCreate := os.Create(opts.To)
		if errCreate != nil {
			panic("can't create destination file")
		}
		bw = bufio.NewWriter(f)
	}
	rw := bufio.NewReadWriter(br, bw)
	err = Copy(rw, opts)
	if err != nil {
		panic("couldn't write")
	}

}
