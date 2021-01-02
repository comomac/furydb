package furydb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"path"
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

	// sanity check find if table exists
	_, table := c.db.findTable(stmt.TableName)
	if table == nil {
		return nil, ErrTableNotExist
	}
	// result remember table schema
	res.tableSchema = table

	// select to all fields
	if stmt.FieldsAll {
		for _, col := range table.Columns {
			stmt.Fields = append(stmt.Fields, col.Name)
		}
	}
	// result remember columns
	res.columns = stmt.Fields

	if Verbose >= 2 {
		fmt.Printf("stmt: %+v\n", stmt)
	}

	folderpath := path.Join(c.db.Folderpath, table.Name)
	res.rows, err = scanDirRows(folderpath, table.Name, stmt.Fields, nil)
	if err != nil {
		return nil, err
	}

	if Verbose >= 2 {
		fmt.Printf("giving results %+v\n", res)
	}

	return res, nil
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

// Where condisions to match
type Where struct {
	OperatorType OperatorType
	Value        interface{}
}

// OperatorType for where comparision
type OperatorType int

// types of comparisions
const (
	OperatorTypeLessThan int = iota
	OperatorTypeLessThanOrEqual
	OperatorTypeMoreThan
	OperatorTypeMoreThanOrEqual
	OperatorTypeEqual
	OperatorTypeNotEqual
)

// scanDirRows scan all the records in table dir for rows
func scanDirRows(folderpath string, tableName string, columns []string, wheres []*Where) ([]*Row, error) {
	files, err := ioutil.ReadDir(folderpath)
	if err != nil {
		return nil, err
	}

	rows := []*Row{}

	for _, file := range files {
		filepath := path.Join(folderpath, file.Name())
		dat, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Printf("read row fail - %s  Err: ( %+v )\n", filepath, err)
			continue
		}

		// row decode
		row := &Row{}
		dec := gob.NewDecoder(bytes.NewReader(dat))
		err = dec.Decode(row)
		if err != nil {
			fmt.Printf("decode row failed - %s  Err: ( %+v )\n", filepath, err)
			continue
		}

		// todo do where match
		// if wheres != nil {
		// 	for _, where := range wheres {
		// ??? continue
		// 	}
		// }

		// sort and filter column accordly
		resCols := []*Column{}
		for _, colName := range columns {
			for _, resCol := range row.Columns {
				if resCol.Name == colName {
					resCols = append(resCols, resCol)
				}
			}
		}
		if len(resCols) != len(columns) {
			fmt.Printf("invalid result column length - %s   %q vs %+v\n", filepath, resCols, columns)
			continue
		}
		row.Columns = resCols

		rows = append(rows, row)
	}

	if Verbose >= 2 {
		fmt.Printf("scanDirRows -> rows %+v\n", rows)
	}

	return rows, nil
}
