package file

import (
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	HashLenagth = 20
	DataFolder = "data"
)

type File struct {
	name string
	hashes string
	recived int
}

func NewFile(name string, hashes string) File {
	return File{
		name:    name,
		hashes:  hashes,
	}
}

func (f *File) AddPart(b []byte) error {
	if !checkPartHash(b, f.getCurrentHash()) {return errors.New("wrong hash")}

	err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", DataFolder, f.name, f.getCurrentHash()), b, 0644)
	if err != nil { return err }

	f.recived++
	return nil
}

func (f *File) getCurrentHash() string {
	return f.hashes[f.recived*HashLenagth:(f.recived+1)*HashLenagth]
}

