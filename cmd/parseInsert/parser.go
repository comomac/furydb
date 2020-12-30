package main

import (
	"fmt"
	"io"
)

// InsertStatement represents a SQL INSERT statement.
type InsertStatement struct {
	FieldsAll bool     // true if using all field(s)
	Fields    []string // or individual field(s)
	Values    []string
	TableName string
}

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

// Parse parses a SQL INSERT statement.
func (p *Parser) Parse() (*InsertStatement, error) {
	stmt := &InsertStatement{}

	// First token should be a "INSERT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != INSERT {
		return nil, fmt.Errorf("found %q, expected INSERT", lit)
	}

	// Next we should see the "INTO" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != INTO {
		return nil, fmt.Errorf("found %q, expected INTO", lit)
	}

	// Next we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	// if we see VALUES, then we know it is all fields
	if tok, _ := p.scanIgnoreWhitespace(); tok == VALUES {
		p.unscan()
		stmt.FieldsAll = true

	} else if tok == LEFTPAR {
		// loop over all our comma-delimited fields.
		var tok Token
		var lit string
		for {
			// Read a field.
			tok, lit = p.scanIgnoreWhitespace()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected field", lit)
			}
			stmt.Fields = append(stmt.Fields, lit)

			// If the next token is not a comma then break the loop.
			if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
				// p.unscan()
				break
			}
		}
		// last token must be )
		if tok != RIGHTPAR {
			return nil, fmt.Errorf("found %q, expected )", lit)
		}
	} else {
		return nil, fmt.Errorf("found %q, unknown state", lit)
	}

	// must be VALUES
	if tok, lit := p.scanIgnoreWhitespace(); tok != VALUES {
		return nil, fmt.Errorf("found %q, expected VALUES", lit)
	}

	// must be (
	if tok, lit := p.scanIgnoreWhitespace(); tok != LEFTPAR {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}

	// Next we should loop over all our comma-delimited values.
	for {
		// Read a value.
		tok, lit = p.scanValueIgnoreWhitespace()
		if tok == VALUE {
			stmt.Values = append(stmt.Values, lit)
		}

		// If the next token is not a comma then break the loop.
		if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
			// p.unscan()
			break
		}
	}

	// last token must be )
	if tok != RIGHTPAR {
		return nil, fmt.Errorf("found %q, expected )))", lit)
	}

	// If the next token is not a comma then break the loop.
	if tok, lit = p.scanValue(); tok != SEMICOL {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	// Return the successfully parsed statement.
	return stmt, nil
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
	if Debug {
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
	if Debug {
		fmt.Printf("tok: %+v      lit: %+v\n", tok, lit)
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
