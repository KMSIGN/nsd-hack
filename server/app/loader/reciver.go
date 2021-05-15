package loader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"nsd-hack/server/app/file"
)

const PartSize = 8 * 1024 * 1024

func SrvFileLoader(filename string, hashes string) (int, error) {
	pt, err := getFreePort()
	if err != nil { return 0, err }
	addr := fmt.Sprintf(":%d", pt)
	log.Printf(addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}

	go func(){
		conn, _ := listener.Accept()
		handle(conn, listener, filename, hashes)
	}()
	return pt, nil
}

func handle(conn net.Conn, listener net.Listener, name string, hashes string) error {
	defer conn.Close()
	defer listener.Close()

	fl := file.NewFile(name, hashes)
	fd := file.NewDownloader(&fl)

	for {
		w := bufio.NewWriter(conn)

		curNo := fd.GetNeededPart()
		_, err := w.WriteString(fmt.Sprintf("%d\n", curNo))
		if err != nil { return err}

		r := bufio.NewReader(conn)
		bts := make([]byte, PartSize)

		n, err := r.Read(bts)
		if err == io.EOF || n != 0 {
			w.WriteString("end\n")
			return nil
		}

		err = fd.AddPart(bts, curNo)
		if err != nil { return err }
		return nil
	}
}
