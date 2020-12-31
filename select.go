package furydb

import (
	"fmt"
	"strings"
)

// querySelect executes a SQL SELECGT statement
func (c *FuryConn) querySelect(query string) (*results, error) {
	res := &results{}

	parser := NewParser(strings.NewReader(query))
	stmt, err := parser.parseSelect()
	if err != nil {
		return nil, err
	}
	fmt.Printf("stmt %+v\n", stmt)

	return res, fmt.Errorf("not implemented")
}

// SelectStatement represents a SQL SELECT statement.
type SelectStatement struct {
	FieldsAll bool     // true if using all field(s)
	Fields    []string // or individual field(s)
	TableName string
}

// parseSelect parses a SQL SELECT statement
func (p *Parser) parseSelect() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
	// if we see *, then we know it is all fields
	if tok, lit := p.scanIgnoreWhitespace(); tok == ASTERISK {
		stmt.FieldsAll = true

	} else if tok == LEFTPAR {
		// loop over all our comma-delimited fields
		for {
			// Read a field.
			tok, lit = p.scanIgnoreWhitespace()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected field", lit)
			}
			stmt.Fields = append(stmt.Fields, lit)

			// If the next token is not a comma then break the loop.
			if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
				break
			}
		}

		// last token must be )
		if tok != RIGHTPAR {
			return nil, fmt.Errorf("found %q, expected )))", lit)
		}
	}

	// Next token should be a "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	// Next we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	// If the next token is not a ; then break the loop.
	if tok, lit = p.scan(); tok != SEMICOL {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	// Return the successfully parsed statement.
	return stmt, nil
}
