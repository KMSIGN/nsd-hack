package filehandler

import (
	"fmt"
	"io"
	"net"
)

func (rcv *FileReciver) CryptRecive(conn net.Conn) (err error) {
	buf := make([]byte, partSize)

	for {

		prtNo := rcv.GetNeededPart()
		if prtNo == -1 {
			fmt.Fprintf(conn, "end\n")
			break
		}
		fmt.Fprintf(conn, "%d\n", prtNo)

		_, err := io.ReadFull(conn, buf)
		if err != nil {
			return err
		}

		//fmt.Printf("enc start:\t %v \n", buf[:15])
		//fmt.Printf("enc end:  \t %v \n", buf[len(buf)-15:])

		err = rcv.AddPart(buf, prtNo)
		if err != nil {
			return err
		}

	}

	return nil
}
