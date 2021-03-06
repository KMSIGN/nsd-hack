package loader

import (
	"fmt"
	"io"
	"net"

	"github.com/KMSIGN/nsd-hack/server/app/file"
)

const PartSize = 8 * 1024 * 1024

func SrvFileLoader(filename string, hashes string, lastPartSize int) (int, error) {
	pt, err := getFreePort()
	if err != nil {
		return 0, err
	}
	addr := fmt.Sprintf(":%d", pt)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, err
	}

	go func() {
		conn, _ := listener.Accept()
		handleRecive(conn, listener, filename, hashes, lastPartSize)
	}()
	return pt, nil
}

func handleRecive(conn net.Conn, listener net.Listener, name string, hashes string, lastPartSize int) error {
	defer conn.Close()
	defer listener.Close()

	fl := file.NewFile(name, hashes, lastPartSize)
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

		n, err := io.ReadFull(conn, buf)
		if err != nil || n == 0 {
			return err
		}

		//fmt.Printf("start:\t %v \n", buf[:15])
		//fmt.Printf("end:  \t %v \n", buf[len(buf)-15:])

		fd.AddPart(buf, curNo)

	}
	return nil
}
