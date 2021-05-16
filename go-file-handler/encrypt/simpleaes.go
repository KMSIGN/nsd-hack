// Make encryption in Go easy
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Aes struct {
	enc, dec cipher.BlockMode
	aes      cipher.Block
	size     int
	key      []byte
	iv       []byte
}

// Creates a new encryption/decryption object
// with a given key of a given size
// (16, 24 or 32 for AES-128, AES-192 and AES-256 respectively,
// as per http://golang.org/pkg/crypto/aes/#NewCipher)
//
// The key will be padded to the given size if needed.
// An IV is created as a series of NULL bytes of necessary length
// when there is no iv string passed as 3rd value to function.
func New(size int, key string, iv []byte) (*Aes, error) {
	padded := make([]byte, size)
	copy(padded, []byte(key))

	if iv == nil {
		iv = make([]byte, size)
		io.ReadFull(rand.Reader, iv)
	}

	aes, err := aes.NewCipher(padded)
	if err != nil {
		return nil, err
	}
	enc := cipher.NewCBCEncrypter(aes, iv)
	dec := cipher.NewCBCDecrypter(aes, iv)
	return &Aes{enc, dec, aes, size, []byte(key), iv}, nil
}

func FromFile(path string) (*Aes, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	contentstr := string(content)
	params := strings.Split(contentstr, ":")
	size, err := strconv.Atoi(params[0])
	if err != nil {
		return nil, err
	}

	return New(size, string(params[1]), []byte(params[2]))
}

func (me *Aes) Save(path string) error {
	content := fmt.Sprintf("%d:%s:%s", me.size, me.key, me.iv)

	err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (me *Aes) padSlice(src []byte) []byte {
	// src must be a multiple of block size
	bs := me.enc.BlockSize()
	mult := int((len(src) / bs) + 1)
	leng := bs * mult

	src_padded := make([]byte, leng)
	copy(src_padded, src)
	return src_padded
}

// Encrypt a slice of bytes, producing a new, freshly allocated slice
//
// Source will be padded with null bytes if necessary
func (me *Aes) Encrypt(src []byte) []byte {
	if len(src)%me.enc.BlockSize() != 0 {
		src = me.padSlice(src)
	}
	dst := make([]byte, len(src))
	me.enc.CryptBlocks(dst, src)
	me.enc = cipher.NewCBCEncrypter(me.aes, me.iv)
	return dst
}

// Encrypt blocks from reader, write results into writer
func (me *Aes) EncryptStream(reader io.Reader, writer io.Writer) error {
	for {
		buf := make([]byte, me.enc.BlockSize())
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF {
				break
			} else if err == io.ErrUnexpectedEOF {
				// nothing
			} else {
				return err
			}
		}
		me.enc.CryptBlocks(buf, buf)
		if _, err = writer.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

// Decrypt a slice of bytes, producing a new, freshly allocated slice
//
// Source will be padded with null bytes if necessary
func (me *Aes) Decrypt(src []byte) []byte {
	if len(src)%me.dec.BlockSize() != 0 {
		src = me.padSlice(src)
	}
	dst := make([]byte, len(src))
	me.dec.CryptBlocks(dst, src)
	return dst
}

// Decrypt blocks from reader, write results into writer
func (me *Aes) DecryptStream(reader io.Reader, writer io.Writer) error {
	buf := make([]byte, me.dec.BlockSize())
	for {
		_, err := io.ReadFull(reader, buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		me.dec.CryptBlocks(buf, buf)
		if _, err = writer.Write(buf); err != nil {
			return err
		}
	}
	return nil
}
