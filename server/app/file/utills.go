package file

import (
	"crypto"
	_ "crypto/sha1"
	"fmt"
)

func checkPartHash(b []byte, h string) bool {
	hash := crypto.SHA1.New()
	hash.Write(b)
	res := hash.Sum(nil)
	return fmt.Sprintf("%x", res) == h
}
