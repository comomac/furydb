package furydb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path"
	"strings"
)

// queryInsert executes a SQL INSERT statement
func (c *FuryConn) queryInsert(query string) (*results, error) {
	res := &results{}

	parser := NewParser(strings.NewReader(query))
	stmt, err := parser.parseInsert()
	if err != nil {
		return nil, err
	}

	// sanity check find if table exists
	_, table := c.db.findTable(stmt.TableName)
	if table == nil {
		return nil, ErrTableNotExist
	}

	// insert to all fields
	if stmt.FieldsAll {
		for _, col := range table.Columns {
			stmt.Fields = append(stmt.Fields, col.Name)
		}
	}

	// sanity check and get formatted columns values
	columns, err := sanityCheckQuery(stmt.Fields, stmt.Values, table)
	if err != nil {
		return nil, err
	}
	// todo probably unnecessary
	// update results
	res.columns = stmt.Fields

	// find row id or generate one
	var pkColName string
	var id string
	for _, cstr := range table.Constraints {
		if cstr.IsPrimaryKey {
			pkColName = cstr.ColumnName
			break
		}
	}
	if pkColName != "" {
		for _, col := range columns {
			if col.Name == pkColName {
				// todo detect id type and pick correctly, instead of using this shortcut
				id = UUIDBinToStr(col.DataUUID)
			}
		}
	}
	if id == "" {
		// todo detect id type correctly, instead of using this shortcut
		id, err = UUIDNewV4()
		if err != nil {
			return nil, err
		}
		idBin, err := UUIDStrToBin(id)
		if err != nil {
			return nil, err
		}
		idCol := &Column{
			Name:        "id",
			Type:        ColumnTypeUUID,
			DataIsNull:  false, // todo false for now
			DataIsValid: true,  // todo valid for now
			DataUUID:    idBin,
		}
		columns = append(columns, idCol)
	}

	// convert to row
	row := &Row{
		TableName: table.Name,
		Columns:   columns,
	}
	// update results
	res.rows = []*Row{row}

	// convert data to bytes
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(row)
	if err != nil {
		return nil, err
	}

	// todo make system support larger than 8k encoded row
	if buf.Len() > 8192 {
		return nil, ErrDataTooBig
	}
	// todo make rows in single file instead of individual files
	filepath := path.Join(c.db.Folderpath, table.Name, id)
	if Verbose >= 3 {
		fmt.Printf("writing row data (%d) to %s", buf.Len(), filepath)
	}
	_, err = writeFile(filepath, row)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// InsertStatement represents a SQL INSERT statement.
type InsertStatement struct {
	FieldsAll bool     // true if using all field(s)
	Fields    []string // or individual field(s)
	Values    []string
	TableName string
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

	// If the next token is not a ; then break the loop.
	if tok, lit = p.scanValue(); tok != SEMICOL {
		return nil, fmt.Errorf("found %q, expected ;", lit)
	}

	// Return the successfully parsed statement.
	return stmt, nil
}
