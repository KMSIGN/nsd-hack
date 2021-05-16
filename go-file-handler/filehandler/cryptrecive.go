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
		_, err = fmt.Fprintf(conn, "%d\n", prtNo)
		if err != nil {
			return err
		}

		if prtNo == *rcv.partCount-1 {
			n, err := io.ReadFull(conn, buf[:rcv.HashUnion.LastPartSize])
			if err != nil || n == 0 {
				return err
			}
		} else {
			n, err := io.ReadFull(conn, buf)
			if err != nil || n == 0 {
				return err
			}
		}

		printBuf("rec enc", buf)

		//decrypted := rcv.encrypter.Decrypt(buf)

		printBuf("rec dec", buf)

		err = rcv.AddPart(buf, prtNo)
		if err != nil {
			return err
		}

	}

	return nil
}
