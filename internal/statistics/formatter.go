package statistics

import (
	"bytes"
	"strconv"
	"time"
	"unicode/utf8"
)

type Formatter interface {
	Format(endpointName string, startAt, endAt time.Time, status int, body []byte) []byte
}

type formatter struct {
	separator rune
	newLine   rune
}

func NewFormatter(separator rune, nl rune) Formatter {
	return &formatter{
		separator: separator,
		newLine:   nl,
	}
}

func (f *formatter) Format(endpointName string, startAt, endAt time.Time, status int, body []byte) []byte {
	sep := make([]byte, utf8.RuneLen(f.separator))
	utf8.EncodeRune(sep, f.separator)
	bnl := make([]byte, utf8.RuneLen(f.newLine))
	utf8.EncodeRune(bnl, f.newLine)
	buffer := bytes.NewBuffer([]byte(endpointName))
	buffer.Write(sep)
	buffer.Write([]byte(endAt.Sub(startAt).String()))
	buffer.Write(sep)
	buffer.WriteString(strconv.Itoa(status))
	buffer.Write(sep)
	buffer.Write(body)
	buffer.Write(bnl)
	return buffer.Bytes()
}
