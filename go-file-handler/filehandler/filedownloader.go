package filehandler

import (
	"errors"
	"os"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

type FileDownloader struct {
	file        *os.File
	encrypter   *encrypt.Aes
	partCount   *int
	hashes      *HashUnion
	neededParts []int
}

func NewDownloader(path string, encrypter *encrypt.Aes, hashes *HashUnion) (*FileDownloader, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	a := make([]int, hashes.PartCount)
	for i := range a {
		a[i] = i
	}

	return &FileDownloader{
		file:        f,
		partCount:   &hashes.PartCount,
		hashes:      hashes,
		neededParts: a,
	}, nil
}

func (fd *FileDownloader) AddPart(b []byte, no int) error {
	curEncHash := fd.hashes.EncHashes[no]
	if !checkHash(b, curEncHash) {
		fd.neededParts = append(fd.neededParts, no)
		return errors.New("wrong encoded hash")
	}

	decrypted := fd.encrypter.Decrypt(b)

	curHash := fd.hashes.Hashes[no]
	if !checkHash(b, curHash) {
		fd.neededParts = append(fd.neededParts, no)
		return errors.New("wrong decrypted hash")
	}
	_, err := fd.file.WriteAt(decrypted, int64(no*partSize))

	if err != nil {
		return err
	}

	return nil
}

func (fd *FileDownloader) GetNeededPart() int {
	if len(fd.neededParts) == 0 {
		return -1
	}
	var x int
	x, fd.neededParts = fd.neededParts[0], fd.neededParts[1:]
	return x
}
