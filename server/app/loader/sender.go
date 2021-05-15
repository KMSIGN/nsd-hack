package loader

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"nsd-hack/server/app/file"
)

func StartUploading(addr string, name string) error {
	if file.CheckFileExists(name) { return errors.New("no such file")}

	// Подключаемся к сокету
	conn, _ := net.Dial("tcp", addr)
	for {

		reader := bufio.NewReader(conn)
		fmt.Print("Text to send: ")
		text, _ := reader.ReadString('\n')
		// Отправляем в socket
		fmt.Fprintf(conn, text + "\n")
		// Прослушиваем ответ
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("Message from server: "+message)
	}


	return nil
}
