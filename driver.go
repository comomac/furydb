package furydb

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"os"
	"strings"
)

// FuryDriver sql driver
type FuryDriver struct {
}

// FuryConn sql connection
type FuryConn struct {
	name string
}

func init() {
	sql.Register("fury", &FuryDriver{})
	fmt.Printf("Drivers=%v\n", sql.Drivers())
}

// Open database
func (d *FuryDriver) Open(name string) (driver.Conn, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return &FuryConn{name: name}, nil
}

// Query the database
func (c *FuryConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	var err error
	var res *results

	str := strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(str, "INSERT") {
		res, err = c.parseQueryInsert(query)
		if err != nil {
			return nil, err
		}
		return res, nil

	} else if strings.HasPrefix(str, "SELECT") {
		res, err = c.parseQueryInsert(query)
		if err != nil {
			return nil, err
		}
		return res, nil

	} else if strings.HasPrefix(str, "UPDATE") {
		res, err = c.parseQueryInsert(query)
		if err != nil {
			return nil, err
		}
		return res, nil

	} else if strings.HasPrefix(str, "DELETE") {
		res, err = c.parseQueryInsert(query)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	return nil, fmt.Errorf("unsupported query")
}

// Begin implements driver.Tx interface, not implemented
func (c *FuryConn) Begin() (_ driver.Tx, err error) {
	return c, fmt.Errorf("Begin method not implemented")
}

// Commit implements driver.Tx interface, not implemented
func (c *FuryConn) Commit() error {
	return fmt.Errorf("Commit method not implemented")
}

// Rollback implements driver.Tx interface, not implemented
func (c *FuryConn) Rollback() error {
	return fmt.Errorf("Rollback method not implemented")
}

// Prepare implements driver.Conn interface, not implemented
func (c *FuryConn) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("Prepare method not implemented")
}

// Close implements driver.Conn interface, not implemented
func (c *FuryConn) Close() error {
	return nil
}
