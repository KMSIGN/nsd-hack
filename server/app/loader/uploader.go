package loader

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
