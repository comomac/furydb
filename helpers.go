package furydb

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var regexUUID = regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}")

// UUIDStrToBin convert uuid string to bytes
func UUIDStrToBin(txt string) (uid [16]byte, err error) {
	dat := strings.ToLower(txt)
	if !regexUUID.MatchString(dat) {
		return uid, ErrInvalidUUID
	}
	dat = strings.ReplaceAll(dat, "-", "")
	b, err := hex.DecodeString(dat)
	if err != nil {
		return uid, err
	}

	// in-place replace
	copy(uid[:], b[:])

	return uid, nil
}

// UUIDBinToStr convert uuid binary to string
func UUIDBinToStr(uid [16]byte) string {
	dst := make([]byte, hex.EncodedLen(len(uid)))
	hex.Encode(dst, uid[:])
	return fmt.Sprintf("%s-%s-%s-%s-%s", dst[0:8], dst[8:12], dst[12:16], dst[16:20], dst[20:32])
}

// UUIDNewV4 generate random uuid string
func UUIDNewV4() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	var b2 [16]byte
	copy(b2[:], b[:])
	return UUIDBinToStr(b2), nil
}
