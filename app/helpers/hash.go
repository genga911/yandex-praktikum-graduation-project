package helpers

import (
	"crypto/md5"
	"encoding/hex"
)

func MakeMD5(in string) string {
	binHash := md5.Sum([]byte(in))
	return hex.EncodeToString(binHash[:])
}
