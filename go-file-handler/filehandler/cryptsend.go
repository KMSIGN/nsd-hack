package filehandler

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/KMSIGN/nsd-hack/go-file-handler/encrypt"
)

func CryptSend(conn net.Conn, file *os.File, key []byte) (err error) {
	reader := bufio.NewReader(file)
	buf := make([]byte, bufferSize)

	cipherText := make([]byte, aes.BlockSize+len(buf))

	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		if n < bufferSize {
			emptyBuffer := [bufferSize]byte{}
			io.ReadFull(bytes.NewReader(buf), emptyBuffer[:n])
			conn.Write(emptyBuffer[:])
		}

		stream.XORKeyStream(cipherText[aes.BlockSize:], buf)

		conn.Write(cipherText)

	}

	return nil
}

func CryptResend(conn net.Conn, file *os.File, encrypter *encrypt.Aes) error {
	buf := make([]byte, bufferSize)
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

		n, err := file.ReadAt(buf, int64(position*bufferSize))
		if err != nil {
			if err != io.EOF {
				return err
			}
		}
		if n == 0 {
			break
		}

		encBuf := encrypter.Encrypt(buf)

		//fmt.Printf("mystr:\t %v \n", encBuf[len(encBuf)-15:])

		conn.Write(encBuf)
	}

	return nil
}
