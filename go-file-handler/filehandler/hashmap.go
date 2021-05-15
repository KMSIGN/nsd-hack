package filehandler

import (
	"bufio"
	"crypto"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

const bufferSize = 8 * 1024 * 1024

func CalcCryptHashMap(f *os.File, enc *encrypt.Aes) (filehash string, hashmap []string, err error) {
	reader := bufio.NewReader(f)
	buf := make([]byte, bufferSize)

	hashes := []string{}

	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		encBuf := enc.Encrypt(buf)

		hash := crypto.SHA1.New()
		hash.Write(encBuf)
		hashes = append(hashes, fmt.Sprintf("%x", hash.Sum(nil)))

	}

	return "", hashes, nil
}

func CalcHashMap(f *os.File) (filehash string, hashmap []string) {
	reader := bufio.NewReader(f)
	buf := make([]byte, bufferSize)
	hash := crypto.SHA1.New()

	hashes := []string{}

	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		hash.Reset()
		hash.Write(buf)
		hashes = append(hashes, fmt.Sprintf("%x", hash.Sum(nil)))

	}

	return "", hashes
}
