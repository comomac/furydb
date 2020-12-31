package furydb

import (
	"fmt"
	"strings"
)

// InsertStatement represents a SQL INSERT statement.
type InsertStatement struct {
	FieldsAll bool     // true if using all field(s)
	Fields    []string // or individual field(s)
	Values    []string
	TableName string
}

// queryInsert executes a SQL INSERT statement
func (c *FuryConn) queryInsert(query string) (*results, error) {
	parser := NewParser(strings.NewReader(query))
	stmt, err := parser.parseInsert()
	if err != nil {
		return nil, err
	}

	// find if table exists
	var table *Table
	for _, tbl := range c.db.Tables {
		if tbl.Name == stmt.TableName {
			table = tbl
		}
	}
	if table == nil {
		return nil, ErrTableNotExist
	}

	if stmt.FieldsAll {
		for _, col := range table.Columns {
			stmt.Fields = append(stmt.Fields, col.Name)
		}
	}

	// sanity check and get formatted columns values
	_, err = sanityCheckQuery(stmt.Fields, stmt.Values, table)
	if err != nil {
		return nil, err
	}

	// todo code data insert

	return nil, nil
}

// parseInsert parses a SQL INSERT statement
func (p *Parser) parseInsert() (*InsertStatement, error) {
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
