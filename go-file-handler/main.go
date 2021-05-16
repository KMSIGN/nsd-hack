package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
	"github.com/KMSIGN/nsd-hack/go-file-handler/filehandler"
)

func main() {
	encrypter, err := encrypt.New(16, "abcgdjfuthatishg", []byte("aaaaaaaaaaaaaaaa"))
	if err != nil {
		log.Fatal(err)
	}
	send(encrypter)
	recive(encrypter)
}

func send(encrypter *encrypt.Aes) {

	fileupl, err := filehandler.NewSender("./test/input/testbin", encrypter)
	if err != nil {
		log.Fatal(err)
	}

	hashjson, err := json.Marshal(fileupl.HashUnion)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("testbin.json", hashjson, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{
		"hash":         {fileupl.HashUnion.SumHash},
		"hashes":       {strings.Join(fileupl.HashUnion.Hashes, "")},
		"lastPartSize": {fmt.Sprint(fileupl.HashUnion.LastPartSize)},
	}

	//fmt.Printf("%s\n", strings.Join(fileupl.HashUnion.EncHashes, " , "))
	//fmt.Printf("%s\n", strings.Join(fileupl.HashUnion.Hashes, " , "))

	resp, err := http.PostForm("http://localhost:8081/upload", data)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	route := "localhost" + string(body)

	conn, err := net.Dial("tcp", route)
	if err != nil {
		log.Fatal(err)
	}

	err = fileupl.CryptSend(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func recive(encrypter *encrypt.Aes) {

	hashjson, err := ioutil.ReadFile("testbin.json")
	if err != nil {
		log.Fatal(err)
	}

	hashunion := &filehandler.HashUnion{}
	err = json.Unmarshal(hashjson, hashunion)
	if err != nil {
		log.Fatal(err)
	}
	filercv, err := filehandler.NewReciver("./test/output/testbin", encrypter, hashunion)
	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{
		"hash": {filercv.HashUnion.SumHash},
	}

	//fmt.Printf("%s\n", strings.Join(fileupl.HashUnion.EncHashes, " , "))
	//fmt.Printf("%s\n", strings.Join(fileupl.HashUnion.Hashes, " , "))

	resp, err := http.PostForm("http://localhost:8081/download", data)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	route := "localhost" + string(body)

	conn, err := net.Dial("tcp", route)
	if err != nil {
		log.Fatal(err)
	}

	err = filercv.CryptRecive(conn)
	if err != nil {
		log.Fatal(err)
	}
}
