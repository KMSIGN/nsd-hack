package filehandler

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"net"
	"os"
)

func CryptRecive(conn net.Conn, file *os.File, hashes []string, key []byte) (err error) {
	writer := bufio.NewWriter(file)
	buf := make([]byte, bufferSize)

	cipherText := make([]byte, aes.BlockSize+len(buf))

	iv := []byte{}
	stream := (cipher.Stream)(nil)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	for {

		if err != nil {
			return err
		}

		if iv == nil {
			iv = cipherText[:aes.BlockSize]
			cipherText = cipherText[aes.BlockSize:]

			stream = cipher.NewCFBDecrypter(block, iv)
		}

		conn.Read(buf)

		stream.XORKeyStream(cipherText[aes.BlockSize:], buf)

		_, err = writer.Write(buf)
		if err != nil {
			return err
		}
	}

	return nil
}
