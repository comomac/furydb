package furydb

import (
	"strings"
)

// parseQueryInsert parse query
func (c *FuryConn) parseQueryInsert(query string) (*results, error) {
	// e.g.
	// INSERT INTO tableName (col1, col2, ...)
	// VALUES (val1, val2, ...);

	query = strings.TrimSpace(query)

	res := &results{}
	// return res, fmt.Errorf("not implemented")
	return res, nil
}
