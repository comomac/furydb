package furydb

import (
	"fmt"
	"strings"
)

// queryTableCreate executes a SQL CREATE TABLE statement
func (c *FuryConn) queryTableCreate(query string) (*results, error) {
	parser := NewParser(strings.NewReader(query))
	_, err := parser.parseTableCreate()
	if err != nil {
		return nil, err
	}

	// todo finish me

	return nil, fmt.Errorf("implement me")
}

// TableCreateStatement represents a SQL CREATE TABLE statement.
type TableCreateStatement struct {
	Columns   []string // of individual column
	Types     []string // of individual column type
	TableName string
}

// parseTableCreate parses a SQL TABLE CREATE statement
func (p *Parser) parseTableCreate() (*TableCreateStatement, error) {
	stmt := &TableCreateStatement{}

	// First token should be a "CREATE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != CREATE {
		return nil, fmt.Errorf("found %q, expected CREATE", lit)
	}

	// Next we should see the "TABLE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != TABLE {
		return nil, fmt.Errorf("found %q, expected TABLE", lit)
	}

	// Next we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table_name", lit)
	}
	stmt.TableName = lit

	// Next we should see the "(" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != LEFTPAR {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}

	// loop over all our comma-delimited column.
	for {
		// Read column.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected column_name", lit)
		}
		stmt.Columns = append(stmt.Columns, lit)

		// Read column type
		tok, lit = p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected column_name", lit)
		}
		stmt.Types = append(stmt.Types, lit)

		// loop over all semi-space seperated column conditions
		// todo

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

	// If the next token is not a comma then break the loop.
	if tok, lit = p.scanValue(); tok != SEMICOL {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	// Return the successfully parsed statement.
	return stmt, nil
}
