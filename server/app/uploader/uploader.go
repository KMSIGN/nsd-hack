package uploader

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"nsd-hack/server/app/file"
)

func SrvFileLoader(filename string, hashes string) (int, error) {
	pt, err := getFreePort()
	if err != nil { return 0, err }
	addr := fmt.Sprintf(":%d", pt)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	conn, err := listener.Accept()
	if err != nil {	return 0, err }
	go handle(conn, filename, hashes)
	return pt, nil
}

func handle(conn net.Conn, name string, hashes string) error {
	defer conn.Close()

	r := bufio.NewReader(conn)
	scanr := bufio.NewScanner(r)

	fl := file.NewFile(name, hashes)

	for {
		scanned := scanr.Scan()
		if !scanned {
			if err := scanr.Err(); err != nil {
				log.Printf("%v(%v)", err, conn.RemoteAddr())
				return err
			}
			break
		}
		err := fl.AddPart(scanr.Bytes())
		if err != nil {	continue }
	}
	return nil
}
