package main

import "C"

import (
	_ "crypto/sha1"
	"strings"

	"github.com/KMSIGN/nsd-hack/go-file-handler/filehandler"
)

//export CalcHashMap
func CalcHashMap(path *C.char) *C.char {
	goPath := C.GoString(path)

	hashes := filehandler.CalcHashMap(goPath)

	return C.CString(strings.Join(hashes, ""))
}

func main() {
	print(C.GoString(CalcHashMap(C.CString("/home/royalcat/neiro.7z"))))
}
