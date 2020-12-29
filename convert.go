package furydb

import (
	"database/sql"
	"errors"
	"fmt"
)

// To help implement NullBytes and NullUUID
// copied some source code from https://golang.org/src/database/sql/sql.go

// convertAssign is the same as convertAssignRows, but without the optional
// rows argument.
func convertAssign(dest, src interface{}) error {
	return convertAssignRows(dest, src, nil)
}

// convertAssignRows copies to dest the value in src, converting it if possible.
// An error is returned if the copy would result in loss of information.
// dest should be a pointer type. If rows is passed in, the rows will
// be used as the parent for any cursor values converted from a
// driver.Rows to a *Rows.
func convertAssignRows(dest, src interface{}, rows *sql.Rows) error {
	// Common cases, without reflect.
	switch s := src.(type) {
	case []byte:
		switch d := dest.(type) {
		case *string:
			if d == nil {
				return errNilPtr
			}
			*d = string(s)
			return nil

		case *interface{}:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil

		case *[]byte:
			if d == nil {
				return errNilPtr
			}
			*d = cloneBytes(s)
			return nil

		case *RawBytes:
			if d == nil {
				return errNilPtr
			}
			*d = s
			return nil
		}
	}

	return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", src, dest)
}

// RawBytes is a byte slice that holds a reference to memory owned by
// the database itself. After a Scan into a RawBytes, the slice is only
// valid until the next call to Next, Scan, or Close.
type RawBytes []byte

var errNilPtr = errors.New("destination pointer is nil") // embedded in descriptive error

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
