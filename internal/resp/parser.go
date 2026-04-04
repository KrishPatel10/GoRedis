package resp

import (
	"bufio"
	"io"
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
		return Value{}, nil
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
	return Value{}, nil
}

func (p *Parser) ParseBulkString() (Value, error) {
	return Value{}, nil
}
