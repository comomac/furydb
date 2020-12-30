package furydb

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
)

// version of furydb
const (
	VersionMajor int = 0 // database schema change
	VersionMinor int = 1 // bug fixes

	Verbose bool = false // use to aid debug
)

// various errors
var (
	ErrTableNotExist      = fmt.Errorf("no such table")
	ErrColumnNotExist     = fmt.Errorf("no such column")
	ErrFieldValueNotMatch = fmt.Errorf("column and values not match")
	ErrValueTypeNotBool   = fmt.Errorf("value type not bool")
	ErrValueTypeNotInt    = fmt.Errorf("value type not int")
	ErrValueTypeNotFloat  = fmt.Errorf("value type not float")
	ErrValueTypeNotString = fmt.Errorf("value type not string")
	ErrValueTypeNotTime   = fmt.Errorf("value type not time")
	ErrValueTypeNotBytes  = fmt.Errorf("value type not bytes")
	ErrValueTypeNotUUID   = fmt.Errorf("value type not uuid")
	ErrValueTypeNotMatch  = fmt.Errorf("value type type not match")
	ErrColumnNotNullable  = fmt.Errorf("column not nullable")
	ErrUnknownColumnType  = fmt.Errorf("unknown column type")
	ErrInvalidUUID        = fmt.Errorf("invalid uuid")
)

// Create new blank database
func Create(folderpath string, name string) (*Database, error) {
	db := &Database{
		Folderpath:   folderpath,
		Name:         name,
		VersionMajor: VersionMajor,
		VersionMinor: VersionMinor,
	}

	return db, nil
}

// Load existing database
func Load(folderpath string) (*Database, error) {
	pathSchema := folderpath + "/schema"

	// file open
	f, err := os.Open(pathSchema)
	if err != nil {
		return nil, err
	}
	// file read
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	// file decode
	db := Database{}
	dec := gob.NewDecoder(bytes.NewReader(data))
	err = dec.Decode(&db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// Save the database to filesystem
func (db *Database) Save(folderpath ...string) error {
	// update to new folderpath
	if len(folderpath) > 0 && folderpath[0] != "" && folderpath[0] != db.Folderpath {
		db.Folderpath = folderpath[0]
	}

	// create dir if not exist
	_, err := os.Stat(db.Folderpath)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(db.Folderpath, 0755)
	} else if err != nil {
		return err
	}

	// save schema. database, table, column
	pathSchema := db.Folderpath + "/schema"
	size, err := writeFile(pathSchema, db)
	if err != nil {
		return err
	}

	if Verbose {
		fmt.Printf("schema %s   size: %d bytes\n", db.Name, size)
	}

	return nil
}

// Close the database
func (db *Database) Close() error {

	return nil
}

// writeFile without retyping lots of code, returns written size
func writeFile(folderpath string, dat interface{}) (int, error) {
	// convert data to bytes
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(dat)
	if err != nil {
		return 0, err
	}

	// file open
	ptr, err := os.Create(folderpath)
	if err != nil {
		return 0, err
	}
	// file write
	size, err := ptr.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}
	// file close
	err = ptr.Close()
	if err != nil {
		return 0, err
	}

	return size, err
}
