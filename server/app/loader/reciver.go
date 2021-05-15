package loader

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/KMSIGN/nsd-hack/server/app/file"
)

const PartSize = 8 * 1024 * 1024

func SrvFileLoader(filename string, hashes string) (int, error) {
	pt, err := getFreePort()
	if err != nil {
		return 0, err
	}
	addr := fmt.Sprintf(":%d", pt)
	log.Printf(addr)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}

	go func() {
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
	buf := make([]byte, PartSize)

	for {

		curNo := fd.GetNeededPart()
		if curNo == -1 {
			_, err := fmt.Fprintf(conn, "end\n")
			if err != nil {
				return err
			}
			break
		}
		_, err := fmt.Fprintf(conn, "%d\n", curNo)
		if err != nil {
			return err
		}

		io.ReadFull(conn, buf)

		//fmt.Printf("mystr:\t %v \n", buf[len(buf)-15:])

		err = fd.AddPart(buf, curNo)
		println(err)
		if err != nil {
			return err
		}

	}
	return nil
}
