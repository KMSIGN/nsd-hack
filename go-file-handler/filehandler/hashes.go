package filehandler

import (
	"bufio"
	"crypto"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

type HashUnion struct {
	PartCount    int      `json:"partCount"`
	LastPartSize int      `json:"lastPartSize"`
	Hashes       []string `json:"hashes"`
	EncHashes    []string `json:"encHashes"`
	SumHash      string   `json:"sumHash"`
}

func NewHashUnionFromFile(file *os.File, encrypter *encrypt.Aes) (*HashUnion, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	hashCount := int(math.Ceil(float64(info.Size()) / float64(partSize)))
	reader := bufio.NewReader(file)
	buf := make([]byte, partSize)

	hash := crypto.SHA1.New()
	encHash := crypto.SHA1.New()
	sumHash := crypto.SHA1.New()

	encHashes := make([]string, hashCount)
	hashes := make([]string, hashCount)

	partIndex := 0
	lastPartSize := 0

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		lastPartSize = n

		sumHash.Write(buf)

		hash.Write(buf)
		hashes[partIndex] = fmt.Sprintf("%x", hash.Sum(nil))
		hash.Reset()

		encBuf := encrypter.Encrypt(buf)

		encHash.Write(encBuf)
		encHashes[partIndex] = fmt.Sprintf("%x", encHash.Sum(nil))
		encHash.Reset()

		partIndex++

	}
	file.Seek(0, 0)

	return &HashUnion{
		PartCount:    hashCount,
		Hashes:       hashes,
		EncHashes:    encHashes,
		LastPartSize: lastPartSize,
		SumHash:      fmt.Sprintf("%x", (sumHash.Sum(nil))),
	}, nil
}
