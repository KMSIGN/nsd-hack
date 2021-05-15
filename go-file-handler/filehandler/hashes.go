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
	PartCount int      `json:"partCount"`
	Hashes    []string `json:"hashes"`
	EncHashes []string `json:"encHashes"`
	SumHash   string   `json:"sumHash"`
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

	encHashes := []string{}
	hashes := []string{}

	for {
		_, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}

		sumHash.Write(buf)

		hash.Write(buf)
		hashes = append(hashes, fmt.Sprintf("%x", encHash.Sum(nil)))
		encHash.Reset()

		encBuf := encrypter.Encrypt(buf)

		encHash.Write(encBuf)
		encHashes = append(hashes, fmt.Sprintf("%x", encHash.Sum(nil)))
		encHash.Reset()

	}
	file.Seek(0, 0)

	return &HashUnion{
		PartCount: hashCount,
		Hashes:    hashes,
		EncHashes: encHashes,
		SumHash:   fmt.Sprintf("%x", (sumHash.Sum(nil))),
	}, nil
}
