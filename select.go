package furydb

import "fmt"

// querySelect executes a SQL SELECGT statement
func (c *FuryConn) querySelect(query string) (*results, error) {
	res := &results{}
	return res, fmt.Errorf("not implemented")
}
