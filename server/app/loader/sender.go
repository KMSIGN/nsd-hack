package loader

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strconv"

	"github.com/KMSIGN/nsd-hack/server/app/file"
)

func StartUploading(addr string, name string) error {
	if file.CheckFileExists(name) { return errors.New("no such file")}

	fl := file.GetFileByName(name)
	fu := file.NewUploader(fl)

	go func () {
		conn, _ := net.Dial("tcp", addr)
		for {
			reader := bufio.NewReader(conn)
			s, err := reader.ReadString('\n')
			if err != nil {
				log.Println("Something go wrong: ", err)
				return
			}

			if s == "end" {	return }

			no, err := strconv.Atoi(s)
			if err != nil {
				log.Println("Bad string to int conv: ", err)
				return
			}

			dat, err := fu.GetPart(no)
			if err != nil {
				log.Println("Getting part error: ", err)
				return
			}

			_, err = conn.Write(dat)
			if err != nil {
				log.Println("Writing to conn error: ", err)
				return
			}
		}
	}()
	return nil
}
