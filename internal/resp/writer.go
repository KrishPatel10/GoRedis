package resp

import (
	"bufio"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
	}
}

func (w *Writer) WriteSimpleString(s string) error {
	return nil
}

func (w *Writer) WriteNull() error {
	_, err := w.writer.Write([]byte("$-1\r\n"))
	return err
}
