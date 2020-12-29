package furydb

import (
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

// Note:
// Table rows are stored separately to make things easier, it may change in the future,
// but for now, the data will just sits in the filesystem not to save memory.
// So data recall or search are all done base on the hard drive, it will be slow on disk hdd,
// but should be quick enough base on SSD and OS caching

// Database holds schema of entire database in self contained unit
type Database struct {
	Name         string // is just a placeholder
	Folderpath   string
	Tables       []*Table
	VersionMajor int
	VersionMinor int
}

// Table holds schema of individual table
type Table struct {
	Name        string
	Columns     []*Column
	Constraints []*Constraint
}

// Constraint holds table column constraint
type Constraint struct {
	Name              string     // name of constraint
	ColumnName        string     // column name for
	Type              ColumnType // what is the type of column
	IsPrimaryKey      bool       // is column primary key
	IsUnique          bool       // is column unique
	IsNotNull         bool       // is column not null
	IsForeignKey      bool       // is this a foreign key?
	ForeignTable      string     // foreign key table
	ForeignColumn     string     // foreign key column
	UseDefaultData    bool       // does it have default value
	DefaultDataBool   bool       // default value in type bool
	DefaultDataInt    int64      // default value in type int64
	DefaultDataFloat  float64    // default value in type float64
	DefaultDataString string     // default value in type string
	DefaultDataBytes  []byte     // default value in type []byte
	DefaultDataTime   string     // default value to use, e.g. now()
	DefaultDataUUID   string     // default value in use, e.g. gen_uuid_v4()
}

// ColumnType that dictate data type that the column value holds
type ColumnType int

// various column types
const (
	ColumnTypeBool   ColumnType = 1
	ColumnTypeInt    ColumnType = 2
	ColumnTypeFloat  ColumnType = 3
	ColumnTypeString ColumnType = 4
	ColumnTypeBytes  ColumnType = 5
	ColumnTypeTime   ColumnType = 6
	ColumnTypeUUID   ColumnType = 7
)

// Column holds schema of individual column, also can be use to hold data
type Column struct {
	Name string     // name of the column
	Type ColumnType // column data type

	// anything below is used for holding data
	DataBool   bool      // value in type bool
	DataInt    int       // value in type int
	DataFloat  float64   // value in type float64
	DataString string    // value in type string
	DataBytes  []byte    // value in type []byte
	DataTime   time.Time // value in type time.Time
	DataUUID   [16]byte  // value in type uuid
}

// Row holds a single row of table data
type Row struct {
	TableName string    // name of the table row refers to
	Data      []*Column // holds column data
	Deleted   bool      // if deleted, will be skipped during scan
}

// results implements driver.Rows
type results struct {
	rows    []*Row
	file    *os.File
	reader  *csv.Reader
	columns []string
}

// Close implements driver.Rows
func (r *results) Close() error {
	return fmt.Errorf("not implemented")
}

// Columns implements driver.Rows
func (r *results) Columns() []string {
	return r.columns
}

// Next implements driver.Rows
func (r *results) Next(dest []driver.Value) error {
	d, err := r.reader.Read()
	if err != nil {
		return err
	}
	for i := 0; i != len(r.columns); i++ {
		dest[i] = driver.Value(d[i])
	}
	return nil
}
