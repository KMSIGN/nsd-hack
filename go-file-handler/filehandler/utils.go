package filehandler

import (
	"crypto"
	"fmt"
)

const partSize = 8 * 1024 * 1024

func checkHash(b []byte, h string) bool {
	hash := crypto.SHA1.New()
	hash.Write(b)
	hashstring := fmt.Sprintf("%x", hash.Sum(nil))
	return hashstring == h
}
