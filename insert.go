package furydb

import "strings"

// parseQueryInsert parse query
func (c *FuryConn) parseQueryInsert(query string) (*results, error) {
	var err error
	// e.g.
	// INSERT INTO tableName (col1, col2, ...)
	// VALUES (val1, val2, ...);

	var sb strings.Builder

	for _, r := range query {
		_, err = sb.WriteRune(r)
		if err != nil {
			return nil, err
		}

		if sb.String() == "INSERT" {
			continue
		}
	}

	res := &results{}
	// return res, fmt.Errorf("not implemented")
	return res, nil
}
