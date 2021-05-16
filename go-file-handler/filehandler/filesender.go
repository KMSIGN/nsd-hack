package filehandler

import (
	"os"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

type FileSender struct {
	file        *os.File
	encrypter   *encrypt.Aes
	partCount   *int
	HashUnion   *HashUnion
	neededParts []int
}

func NewSender(path string, encrypter *encrypt.Aes) (*FileSender, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	hashunion, err := NewHashUnionFromFile(f, encrypter)
	if err != nil {
		return nil, err
	}

	a := make([]int, hashunion.PartCount)
	for i := range a {
		a[i] = i
	}

	return &FileSender{
		file:        f,
		partCount:   &hashunion.PartCount,
		encrypter:   encrypter,
		HashUnion:   hashunion,
		neededParts: a,
	}, nil
}
