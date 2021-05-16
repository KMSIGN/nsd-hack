package filehandler

import (
	"errors"
	"log"
	"os"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

type FileReciver struct {
	file        *os.File
	encrypter   *encrypt.Aes
	partCount   *int
	HashUnion   *HashUnion
	neededParts []int
}

func NewReciver(path string, encrypter *encrypt.Aes, hashes *HashUnion) (*FileReciver, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	a := make([]int, hashes.PartCount)
	for i := range a {
		a[i] = i
	}

	return &FileReciver{
		file:        f,
		partCount:   &hashes.PartCount,
		HashUnion:   hashes,
		encrypter:   encrypter,
		neededParts: a,
	}, nil
}

func (fd *FileReciver) AddPart(b []byte, no int) error {

	curHash := fd.HashUnion.Hashes[no]
	if !checkHash(b, curHash) {
		fd.neededParts = append(fd.neededParts, no)
		log.Println("wrong decrypted hash retrying")
		return errors.New("wrong decrypted hash")

	}

	if no == *fd.partCount-1 {
		_, err := fd.file.WriteAt(b[:fd.HashUnion.LastPartSize], int64(no*partSize))
		if err != nil {
			return err
		}
	} else {
		_, err := fd.file.WriteAt(b, int64(no*partSize))
		if err != nil {
			return err
		}
	}

	return nil
}

func (fd *FileReciver) GetNeededPart() int {
	if len(fd.neededParts) == 0 {
		return -1
	}
	var x int
	x, fd.neededParts = fd.neededParts[0], fd.neededParts[1:]
	return x
}
