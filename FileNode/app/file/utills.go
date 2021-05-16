package file

import (
	"crypto"
	_ "crypto/sha1"
	"fmt"
	"io/ioutil"
)

func checkPartHash(b []byte, h string) bool {
	hash := crypto.SHA1.New()
	hash.Write(b)
	hashstring := fmt.Sprintf("%x", hash.Sum(nil))
	return hashstring == h
}

func GetFilesLen() int {
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/", DataFolder))
	return len(files)
}
