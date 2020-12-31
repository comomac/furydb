package main

import (
	"fmt"
	"testing"

	"github.com/comomac/furydb"
)

// TestUUIDNewV4 generate uuid v4
func TestUUIDNewV4(t *testing.T) {
	uid, err := furydb.UUIDNewV4()
	if err != nil {
		t.Error(err)
	}
	if len(uid) != 36 {
		t.Error(fmt.Errorf("invalid uuid string len"))
	}

	buid, err := furydb.UUIDStrToBin(uid)
	if err != nil {
		t.Error(err)
	}
	uid2 := furydb.UUIDBinToStr(buid)

	if uid != uid2 {
		t.Error(fmt.Errorf("uuid conversion fail"))
	}
}
