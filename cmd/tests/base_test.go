package fury

import (
	"fmt"
	"testing"

	"github.com/comomac/furydb"
)

// TestCreate db
func TestCreate(t *testing.T) {
	db, err := furydb.Create("tmp-db", "testme")
	if err != nil {
		t.Error(err)
	}

	err = db.Save()
	if err != nil {
		t.Error(err)
	}
}

//d TestLoad db
func TestLoad(t *testing.T) {
	db, err := furydb.Load("tmp-db")
	if err != nil {
		t.Error(err)
	}

	if db.Name != "testme" {
		t.Error(fmt.Errorf("name mismatch"))
	}
}
