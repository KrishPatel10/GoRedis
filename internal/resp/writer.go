package resp

import (
	"bufio"
	"io"
	"strconv"
)

type Writer struct {
	// Change this from io.Writer to *bufio.Writer so we have access to Flush()
	writer *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
	}
}

func (w *Writer) WriteSimpleString(s string) error {
	w.writer.WriteString("+" + s + "\r\n")
	return w.writer.Flush() // Push it over the TCP network
}

// WriteBulkString is required for the GET command!
func (w *Writer) WriteBulkString(s string) error {
	w.writer.WriteString("$" + strconv.Itoa(len(s)) + "\r\n")
	w.writer.WriteString(s + "\r\n")
	return w.writer.Flush()
}

func (w *Writer) WriteNull() error {
	w.writer.WriteString("$-1\r\n")
	return w.writer.Flush()
}
