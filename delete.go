package furydb

import "fmt"

// queryDelete executes a SQL DELETE statement
func (c *FuryConn) queryDelete(query string) (*results, error) {
	res := &results{}
	return res, fmt.Errorf("not implemented")
}
