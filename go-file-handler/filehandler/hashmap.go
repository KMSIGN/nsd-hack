package filehandler

import (
	"bufio"
	"crypto"
	"fmt"
	"io"
	"log"
	"os"
)

const bufferSize = 8 * 1024 * 1024

func CalcHashMap(path string) []string {
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

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

	return hashes
}
