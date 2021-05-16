package filehandler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func (snd *FileSender) CryptSend(conn net.Conn) error {
	buf := make([]byte, partSize)
	connReader := bufio.NewReader(conn)

	for {
		message, err := connReader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.HasPrefix(message, "end") {
			break
		}

		position, err := strconv.Atoi(message[:len(message)-1])
		if err != nil {
			return err
		}

		_, err = snd.file.ReadAt(buf, int64(position*partSize))
		if err != nil {
			if err != io.EOF {
				return err
			}
		}

		fmt.Printf("send dec start:\t %v \n", buf[:15])
		fmt.Printf("send dec end:  \t %v \n", buf[len(buf)-15:])

		encBuf := snd.encrypter.Encrypt(buf)

		//fmt.Printf("send enc start:\t %v \n", encBuf[:15])
		//fmt.Printf("send enc end:  \t %v \n", encBuf[len(encBuf)-15:])

		conn.Write(encBuf)
	}

	return nil
}
