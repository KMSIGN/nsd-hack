package api

import (
	"crypto"
	_ "crypto/sha1"
	"fmt"
)

func nameHashing(x string) string {
	hash := crypto.SHA1.New()
	hash.Write([]byte(x))
	res := hash.Sum(nil)
	return fmt.Sprintf("%x", res)
}
