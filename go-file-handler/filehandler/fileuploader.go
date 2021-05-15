package filehandler

import (
	"os"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

type FileUploader struct {
	file        *os.File
	encrypter   *encrypt.Aes
	partCount   *int
	hashes      *HashUnion
	neededParts []int
}

func NewUploader(path string, hashes *HashUnion, encrypter *encrypt.Aes) (*FileUploader, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	a := make([]int, hashes.PartCount)
	for i := range a {
		a[i] = i
	}

	return &FileUploader{
		file:        f,
		partCount:   &hashes.PartCount,
		encrypter:   encrypter,
		hashes:      hashes,
		neededParts: a,
	}, nil
}
