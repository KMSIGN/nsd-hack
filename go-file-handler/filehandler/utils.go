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

func printBuf(comment string, b []byte) {
	fmt.Printf("%s start:\t %v \n", comment, b[:16])
	fmt.Printf("%s mid:  \t %v \n", comment, b[len(b)/2-8:len(b)/2+8])
	fmt.Printf("%s end:  \t %v \n", comment, b[len(b)-16:])
}
