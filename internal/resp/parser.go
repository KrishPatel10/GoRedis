package resp

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Value struct {
	Typ   string
	Value string
	Array []Value
}

type Parser struct {
	reader *bufio.Reader
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(r),
	}
}

func (p *Parser) Parse() (Value, error) {
	typ, err := p.reader.ReadByte()

	if err != nil {
		// FIXED: We must return the error here so the server knows the connection dropped!
		return Value{}, err
	}

	switch typ {
	case '*':
		return p.ParseArray()
	case '$':
		return p.ParseBulkString()
	default:
		return Value{}, io.ErrUnexpectedEOF
	}
}

func (p *Parser) ParseArray() (Value, error) {
	// 1. Read up to the first newline
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	// 2. Remove the trailing \r\n
	line = strings.TrimSuffix(line, "\r\n")

	// 3. Convert the length string to an integer
	count, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}

	// 4. Loop count times, parsing each element
	var array []Value
	for i := 0; i < count; i++ {
		val, err := p.Parse()
		if err != nil {
			return Value{}, err
		}
		array = append(array, val)
	}

	return Value{
		Typ:   "array",
		Array: array,
	}, nil
}

func (p *Parser) ParseBulkString() (Value, error) {
	// 1. Read up to the newline to get the length
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	// 2. Remove trailing \r\n and convert to integer
	line = strings.TrimSuffix(line, "\r\n")
	length, err := strconv.Atoi(line)
	if err != nil {
		return Value{}, err
	}

	// 3. Create a buffer of the exact string length and read the data
	buf := make([]byte, length)
	_, err = io.ReadFull(p.reader, buf)
	if err != nil {
		return Value{}, err
	}

	// 4. THE CRITICAL FIX: Read and discard the trailing \r\n
	p.reader.ReadByte() // consumes \r
	p.reader.ReadByte() // consumes \n

	return Value{
		Typ:   "bulk",
		Value: string(buf),
	}, nil
}
