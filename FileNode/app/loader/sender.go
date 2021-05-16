package loader

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/KMSIGN/nsd-hack/server/app/file"
)

func StartUploading(name string) (int, error) {
	pt, err := getFreePort()
	if err != nil {
		return 0, err
	}

	if !file.CheckFileExists(name) {
		return 0, errors.New("no such file")
	}

	addr := fmt.Sprintf(":%d", pt)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}

	fl := file.GetFileByName(name)
	fu := file.NewUploader(fl)

	go func() {
		conn, _ := listener.Accept()
		handleSend(conn, listener, fu)
	}()
	return pt, nil
}

func handleSend(conn net.Conn, listener net.Listener, fu *file.UploaderFile) error {
	for {
		reader := bufio.NewReader(conn)
		s, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Something go wrong: ", err)
			return err
		}

		if strings.HasPrefix(s, "end") {
			return err
		}

		n, err := strconv.Atoi(s[:len(s)-1])
		if err != nil {
			log.Println("Bad string to int conv: ", err)
			return err
		}

		dat, err := fu.GetPart(n)
		if err != nil {
			log.Println("Getting part error: ", err)
			return err
		}

		//fmt.Printf("start:\t %v \n", dat[:15])
		//fmt.Printf("end:  \t %v \n", dat[len(dat)-15:])

		_, err = conn.Write(dat)
		if err != nil {
			log.Println("Writing to conn error: ", err)
			return err
		}
	}
}
