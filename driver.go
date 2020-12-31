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
	db *Database
}

func init() {
	sql.Register("fury", &FuryDriver{})
	fmt.Printf("Drivers=%v\n", sql.Drivers())
}

// Open database
func (d *FuryDriver) Open(folderPath string) (driver.Conn, error) {
	filePath := folderPath + "/schema"

	// file not exist -> new
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {

	} else if err != nil {
		return nil, err
	}

	// load file
	db, err := Load(folderPath)
	if err != nil {
		return nil, err
	}

	return &FuryConn{db: db}, nil
}

// Query the database
func (c *FuryConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	var err error
	var res *results

	str := strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(str, "INSERT") {
		res, err = c.queryInsert(query)
		if err != nil {
			return nil, err
		}
		return res, nil

	} else if strings.HasPrefix(str, "SELECT") {
		res, err = c.querySelect(query)
		if err != nil {
			return nil, err
		}
		return res, nil

	} else if strings.HasPrefix(str, "UPDATE") {
		res, err = c.queryUpdate(query)
		if err != nil {
			return nil, err
		}
		return res, nil

	} else if strings.HasPrefix(str, "DELETE") {
		res, err = c.queryDelete(query)
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
