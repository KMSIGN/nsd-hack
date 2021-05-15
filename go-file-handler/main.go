package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
	"github.com/KMSIGN/nsd-hack/go-file-handler/filehandler"
)

func main() {
	encrypter, err := encrypt.New(16, "abcgdjfuthatishg")
	if err != nil {
		panic(err)
	}

	file, err := os.Open("/home/royalcat/neiro.7z")
	if err != nil {
		panic(err)
	}
	_, hashmap, err := filehandler.CalcCryptHashMap(file, encrypter)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", strings.Join(hashmap, " , "))
	file.Seek(0, 0)

	data := url.Values{"hash": {"/home/royalcat/neiro.7z"}, "hashes": {strings.Join(hashmap, "")}}
	resp, err := http.PostForm("http://localhost:8080/upload", data)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	route := "localhost" + string(body)

	conn, err := net.Dial("tcp", route)
	if err != nil {
		panic(err)
	}

	err = filehandler.CryptResend(conn, file, encrypter)
	if err != nil {
		panic(err)
	}
}
