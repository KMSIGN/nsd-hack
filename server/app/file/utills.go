package file

import (
	"crypto"
	_ "crypto/sha1"
	"fmt"
)

func checkPartHash(b []byte, h string) bool {
	hash := crypto.SHA1.New()
	hash.Write(b)
	hashstring := fmt.Sprintf("%x", hash.Sum(nil))
	return hashstring == h
}
