package furydb

import (
	"database/sql"
	"database/sql/driver"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
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
	DefaultDataTime   string     // default value to use, e.g. now()
	DefaultDataBytes  []byte     // default value in type []byte
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
	ColumnTypeTime   ColumnType = 5
	ColumnTypeBytes  ColumnType = 6
	ColumnTypeUUID   ColumnType = 7
)

// Column holds schema of individual column, also can be use to hold data
type Column struct {
	Name string     // name of the column
	Type ColumnType // column data type

	// anything below is used for holding data
	DataIsNull bool      // value is null (if column is nullable)
	DataBool   bool      // value in type bool
	DataInt    int64     // value in type int
	DataFloat  float64   // value in type float64
	DataString string    // value in type string
	DataTime   time.Time // value in type time.Time
	DataBytes  []byte    // value in type []byte
	DataUUID   [16]byte  // value in type uuid
}

// Row holds a single row of table data
type Row struct {
	TableName string    // name of the table row refers to
	Columns   []*Column // holds column data
	Deleted   bool      // if deleted, will be skipped during scan
}

// results implements driver.Rows
type results struct {
	tableSchema *Table
	rows        []*Row
	reader      *csv.Reader
	cursor      int // increment after each Next()
	columns     []string
}

// Close implements driver.Rows
func (r *results) Close() error {
	// return fmt.Errorf("not implemented")
	return nil
}

// Columns implements driver.Rows
func (r *results) Columns() []string {
	return r.columns
}

// Next implements driver.Rows
func (r *results) Next(dest []driver.Value) error {
	constraints := r.tableSchema.Constraints

	row := r.rows[r.cursor]
	for i, col := range row.Columns {
		var constraint *Constraint
		for _, cstr := range constraints {
			if cstr.ColumnName == col.Name {
				constraint = cstr
			}
		}

		switch col.Type {
		case ColumnTypeBool:
			if constraint == nil {
				dest[i] = driver.Value(col.DataBool)
			} else {
				dest[i] = sql.NullBool{
					Bool:  col.DataBool,
					Valid: col.DataIsNull,
				}
			}
		case ColumnTypeInt:
			if constraint == nil {
				dest[i] = driver.Value(col.DataInt)
			} else {
				dest[i] = sql.NullInt64{
					Int64: col.DataInt,
					Valid: col.DataIsNull,
				}
			}
		case ColumnTypeFloat:
			if constraint == nil {
				dest[i] = driver.Value(col.DataFloat)
			} else {
				dest[i] = sql.NullFloat64{
					Float64: col.DataFloat,
					Valid:   col.DataIsNull,
				}
			}
		case ColumnTypeString:
			if constraint == nil {
				dest[i] = driver.Value(col.DataString)
			} else {
				dest[i] = sql.NullString{
					String: col.DataString,
					Valid:  col.DataIsNull,
				}
			}
		case ColumnTypeTime:
			if constraint == nil {
				dest[i] = driver.Value(col.DataTime)
			} else {
				dest[i] = sql.NullTime{
					Time:  col.DataTime,
					Valid: col.DataIsNull,
				}
			}
		case ColumnTypeBytes:
			if constraint == nil {
				dest[i] = driver.Value(col.DataBytes)
			} else {
				dest[i] = NullBytes{
					Bytes: col.DataBytes,
					Valid: col.DataIsNull,
				}
			}
		case ColumnTypeUUID:
			if constraint == nil {
				dest[i] = driver.Value(col.DataUUID)
			} else {
				dest[i] = NullUUID{
					UUID:  col.DataUUID,
					Valid: col.DataIsNull,
				}
			}
		default:
			return fmt.Errorf("unsupported column type", col.Type)
		}
	}
	return nil
}

// NullBytes for nullable bytes
type NullBytes struct {
	Bytes []byte
	Valid bool
}

// Scan implements the Scanner interface
func (n *NullBytes) Scan(value interface{}) error {
	if value == nil {
		n.Bytes, n.Valid = []byte{}, false
		return nil
	}
	n.Valid = true
	// return convertAssign(&n.Bytes, value)

	// see if this will work
	dat, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot Scan NullBytes value")
	}
	copy(n.Bytes, dat)
	return nil

}

// Value implements the Valuer interface
func (n *NullBytes) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}

// NullUUID for nullable uuid
type NullUUID struct {
	UUID  [16]byte
	Valid bool
}

var regexUUID = regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}")

// Scan implements the Scanner interface
func (n *NullUUID) Scan(value interface{}) error {
	if value == nil {
		n.UUID, n.Valid = [16]byte{}, false
		return nil
	}
	n.Valid = true
	// return convertAssign(&n.UUID, value)

	// convert hex string (uuid) to byte
	dat, ok := value.(string)
	if !ok {
		return fmt.Errorf("not a string")
	}
	dat = strings.ToLower(dat)
	if !regexUUID.MatchString(dat) {
		return fmt.Errorf("invalid uuid")
	}
	dat = strings.ReplaceAll(dat, "-", "")
	b, err := hex.DecodeString(dat)
	if err != nil {
		return err
	}

	// in-place replace
	copy(n.UUID[:], b[:])
	return nil
}

// Value implements the Valuer interface
func (n *NullUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return "00000000-0000-0000-0000-000000000000", nil
	}

	dst := make([]byte, hex.EncodedLen(len(n.UUID)))
	hex.Encode(dst, n.UUID[:])

	str := fmt.Sprintf("%s-%s-%s-%s-%s", dst[0:7], dst[8:11], dst[12:15], dst[16:19], dst[20:])

	return str, nil
}
