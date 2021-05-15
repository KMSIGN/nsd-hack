package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

const (
	HashLenagth = 20
	DataFolder = "data"
)

type File struct {
	Hash    string `json:"hash"`
	Hashes  string `json:"hashes"`
}

func NewFile(hash string, hashes string) File {
	fl := File{
		Hash:   hash,
		Hashes: hashes,
	}
	t, _ := json.Marshal(fl)
	ioutil.WriteFile(fmt.Sprintf("%s/%s/scheme.json", DataFolder, fl.Hash), t, 0644)
	return fl
}

func (f *File) GetHashByNo(n int) string {
	return f.Hashes[n*HashLenagth:(n+1)*HashLenagth]
}

type DownloaderFile struct {
	file *File
	neededParts []int
}

func NewDownloader(f *File) *DownloaderFile{
	a := make([]int, len(f.Hashes)/HashLenagth)
	for i := range a {
		a[i] = i
	}
	return &DownloaderFile{
		file:        f,
		neededParts: a,
	}
}

func (fd *DownloaderFile) AddPart(b []byte, no int) error {
	curHash := fd.file.GetHashByNo(no)
	if !checkPartHash(b, curHash) {
		fd.neededParts = append(fd.neededParts, no)
		return errors.New("wrong hash")
	}

	err := ioutil.WriteFile(fmt.Sprintf("%s/%s/%s", DataFolder, fd.file.Hash, curHash), b, 0644)
	if err != nil { return err }


	return nil
}

func (fd *DownloaderFile) GetNeededPart() int {
	if len(fd.neededParts) == 0{
		return -1
	}
	var x int
	x, fd.neededParts = fd.neededParts[0], fd.neededParts[1:]
	return x
}

func CheckFileExists(n string) bool {
	files, err := ioutil.ReadDir(fmt.Sprintf("%s/", DataFolder))
	if err != nil {
		log.Println(err)
		return false
	}

	for _, f := range files {
		if n == f.Name() { return true }
	}
	return false
}

func GetFileByName(n string) *File {
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/", DataFolder))

	for _, f := range files {
		if n == f.Name() {
			bt, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s/scheme.json", DataFolder, f.Name()))
			var file File
			err := json.Unmarshal(bt, &file)
			if err != nil { return nil}
			return &file
		}
	}
	return nil
}

type UploaderFile struct {
	file *File
	NeedsToUpload []int
}

func NewUploader(f *File) *UploaderFile {
	a := make([]int, len(f.Hashes)/HashLenagth)
	for i := range a {
		a[i] = i
	}
	return &UploaderFile{
		file: f,
		NeedsToUpload: a,
	}
}

func (fu *UploaderFile) ErrorInUploading(no int){
	fu.NeedsToUpload = append(fu.NeedsToUpload, no)
}

func (fu *UploaderFile) GetPart() ([]byte, int, error) {
	n := fu.getNextPartNo()
	hs := fu.file.GetHashByNo(n)
	bts, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/%s", DataFolder, fu.file.Hash, hs))
	if err != nil { return nil, -1, err}
	return bts, n, nil
}

func (fu *UploaderFile) getNextPartNo() int {
	var x int
	x, fu.NeedsToUpload = fu.NeedsToUpload[0], fu.NeedsToUpload[1:]
	return x
}