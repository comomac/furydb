package furydb

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var regexUUID = regexp.MustCompile("[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}")

func uuidStrToBin(txt string) (uid [16]byte, err error) {
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

func uuidBinToStr(uid [16]byte) string {
	dst := make([]byte, hex.EncodedLen(len(uid)))
	hex.Encode(dst, uid[:])

	str := fmt.Sprintf("%s-%s-%s-%s-%s", dst[0:7], dst[8:11], dst[12:15], dst[16:19], dst[20:])
	return str
}
