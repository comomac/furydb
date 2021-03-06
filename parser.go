package furydb

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scanValue() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.ScanValue()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanValueIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scanValue()
	if tok == WS {
		tok, lit = p.scanValue()
	}
	if Verbose >= 3 {
		fmt.Printf("tok: %+v      lit: %+v\n", tok, lit)
	}
	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	if Verbose >= 3 {
		fmt.Printf("tok: %+v      lit: %+v\n", tok, lit)
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// sanityCheckQuery check the field and value, and return formatted columns
func sanityCheckQuery(fields []string, values []string, table *Table) ([]*Column, error) {
	// result columns with data
	rColumns := []*Column{}

	if Verbose >= 3 {
		fmt.Printf("fields: (%d) %q\n", len(fields), fields)
		fmt.Printf("values: (%d) %q\n", len(values), values)
	}

	// sanity check fields and values length
	if len(fields) != len(values) {
		return nil, ErrFieldValueLengthNotMatch
	}
	// match value type and column data type
	for i, field := range fields {
		// find if column exists
		var column *Column
		for _, col := range table.Columns {
			if col.Name == field {
				// duplicate so we dont mutate the original column
				column = &*col
			}
		}
		if column == nil {
			return nil, ErrColumnNotExist
		}
		var constraint *Constraint
		for _, cstr := range table.Constraints {
			if cstr.ColumnName == field {
				constraint = cstr
			}
		}

		if Verbose >= 4 {
			fmt.Printf("column: %+v\n", column)
		}

		value := values[i]
		switch column.Type {
		case ColumnTypeBool:
			switch strings.ToLower(value) {
			case "true":
				column.DataBool = true
			case "false":
				column.DataBool = false
			case "null":
				if constraint == nil {
					return nil, ErrColumnNotNullable
				}
				column.DataBool = false
				column.DataIsNull = true
			default:
				return nil, ErrValueTypeNotBool
			}
		case ColumnTypeInt:
			if strings.ToLower(value) == "null" {
				if constraint == nil {
					return nil, ErrColumnNotNullable
				}
				column.DataInt = 0
				column.DataIsNull = true
			} else {
				num, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, ErrValueTypeNotInt
				}
				column.DataInt = num
			}
		case ColumnTypeFloat:
			if strings.ToLower(value) == "null" {
				if constraint == nil {
					return nil, ErrColumnNotNullable
				}
				column.DataFloat = 0
				column.DataIsNull = true
			} else {
				num, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, ErrValueTypeNotFloat
				}
				column.DataFloat = num
			}
		case ColumnTypeString:
			// todo, null or 'null' is just treated as null, this could be problematic
			if strings.ToLower(value) == "null" {
				if constraint == nil {
					return nil, ErrColumnNotNullable
				}
				column.DataString = ""
				column.DataIsNull = true
			} else {
				column.DataString = value
			}
		case ColumnTypeTime:
			if strings.ToLower(value) == "null" {
				if constraint == nil {
					return nil, ErrColumnNotNullable
				}
				column.DataTime = time.Time{}
				column.DataIsNull = true
			} else {
				t, err := time.Parse(value, time.RFC3339)
				if err != nil {
					return nil, ErrValueTypeNotTime
				}
				column.DataTime = t
			}
		case ColumnTypeUUID:
			if strings.ToLower(value) == "null" {
				if constraint == nil {
					return nil, ErrColumnNotNullable
				}
				column.DataUUID = [16]byte{}
				column.DataIsNull = true
			} else {
				b, err := UUIDStrToBin(value)
				if err != nil {
					return nil, ErrValueTypeNotUUID
				}
				column.DataUUID = b
			}
		default:
			return nil, ErrUnknownColumnType
		}
		// add columns together and eventually use for return
		rColumns = append(rColumns, column)
	}

	return rColumns, nil
}
