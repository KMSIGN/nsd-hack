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
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "filedu",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "addres",
				Aliases: []string{"a"},
				Usage:   "address of main server (like localhost)",
			},
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "port of main server (like 8081)",
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "Load from/to `FILE`",
			},
			&cli.StringFlag{
				Name:        "hashes",
				DefaultText: "hashes.json",
				Aliases:     []string{"H"},
				Usage:       "Load from/to `FILE`",
			},
			&cli.StringFlag{
				Name:        "key",
				DefaultText: "keys.key",
				Aliases:     []string{"k"},
				Usage:       "Load from/to `FILE`",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "upload",
				Aliases: []string{"c"},
				Usage:   "upload file to server",
				Action: func(c *cli.Context) error {
					filePath := c.String("file")
					addres := c.String("addres")
					port := c.String("port")
					keys := c.String("key")
					if keys == "" {
						keys = "keys.key"
					}
					hashesJson := c.String("hashes")
					if hashesJson == "" {
						hashesJson = "hashes.json"
					}
					send(filePath, keys, hashesJson, addres, port)
					return nil
				},
			},
			{
				Name:    "download",
				Aliases: []string{"a"},
				Usage:   "download file form server",
				Action: func(c *cli.Context) error {
					filePath := c.String("file")
					addres := c.String("addres")
					port := c.String("port")
					keys := c.String("key")
					if keys == "" {
						keys = "keys.key"
					}
					hashesJson := c.String("hashes")
					if hashesJson == "" {
						hashesJson = "hashes.json"
					}
					recive(filePath, keys, hashesJson, addres, port)
					return nil
				},
			},
		}}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func send(filepath string, keypath string, hashpath string, addres string, port string) {

	encrypter, err := encrypt.New(16, "abcgdjfuthatishg", []byte("aaaaaaaaaaaaaaaa"))
	if err != nil {
		log.Fatal(err)
	}

	encrypter.Save(keypath)

	fileupl, err := filehandler.NewSender(filepath, encrypter)
	if err != nil {
		log.Fatal(err)
	}

	hashjson, err := json.Marshal(fileupl.HashUnion)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(hashpath, hashjson, os.ModePerm)
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

	resp, err := http.PostForm(fmt.Sprintf("http://%s:%s/upload", addres, port), data)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	route := addres + string(body)

	conn, err := net.Dial("tcp", route)
	if err != nil {
		log.Fatal(err)
	}

	err = fileupl.CryptSend(conn)
	if err != nil {
		log.Fatal(err)
	}
}

func recive(filepath string, keypath string, hashpath string, addres string, port string) {

	encrypter, err := encrypt.FromFile(keypath)
	if err != nil {
		log.Fatal(err)
	}

	hashjson, err := ioutil.ReadFile(hashpath)
	if err != nil {
		log.Fatal(err)
	}

	hashunion := &filehandler.HashUnion{}
	err = json.Unmarshal(hashjson, hashunion)
	if err != nil {
		log.Fatal(err)
	}
	filercv, err := filehandler.NewReciver(filepath, encrypter, hashunion)
	if err != nil {
		log.Fatal(err)
	}

	data := url.Values{
		"hash": {filercv.HashUnion.SumHash},
	}

	//fmt.Printf("%s\n", strings.Join(fileupl.HashUnion.EncHashes, " , "))
	//fmt.Printf("%s\n", strings.Join(fileupl.HashUnion.Hashes, " , "))

	resp, err := http.PostForm(fmt.Sprintf("http://%s:%s/download", addres, port), data)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(string(body))
	}

	route := addres + string(body)

	conn, err := net.Dial("tcp", route)
	if err != nil {
		log.Fatal(err)
	}

	err = filercv.CryptRecive(conn)
	if err != nil {
		log.Fatal(err)
	}
}
