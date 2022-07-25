package security

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
)

func CreateName(name string) string {
	h := md5.New()
	tmpName :=
		strconv.FormatInt(time.Now().Unix(), 10) +
			name +
			strconv.FormatInt(time.Now().UnixNano()%0x100000, 10)
	h.Write([]byte(tmpName))

	return hex.EncodeToString(h.Sum(nil))
}
