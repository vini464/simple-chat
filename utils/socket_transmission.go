package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

func receiveMessage(conn net.Conn, buffer []byte) error {
	receicedLen := 0
	for receicedLen < len(buffer) {
		readed, err := conn.Read(buffer[receicedLen:]) // escreve as próximas informações sempre na frente
		if err != nil {
			return err
		}
		receicedLen += readed
	}
	return nil
}

func ReceiveHandler(conn net.Conn, received_data chan []byte, wg *sync.WaitGroup) {
	header := make([]byte, 4) // tamanho da informação que virá, sempre será 4 bytes
	for {
		err := receiveMessage(conn, header)
		if err != nil {
			fmt.Println("[error] - error while reading message:", err)
      wg.Done()
			return
		}
		msg_size := binary.BigEndian.Uint32(header)
		data := make([]byte, msg_size)
		err = receiveMessage(conn, data)
		if err != nil {
			fmt.Println("[error] - error while readign message:", err)
      wg.Done()
			return
		}
		received_data <- data
	}
}

func SendHandler(conn net.Conn, msg chan []byte, wg *sync.WaitGroup) {
	for {
		data := <-msg
		size := uint32(len(data))
		header := make([]byte, 4)

		binary.BigEndian.PutUint32(header, size)
		_, err := conn.Write(header)
		if err != nil {
			fmt.Println("[error] - error while sending message:", err)
      wg.Done()
			return
		}
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("[error] - error while sending message:", err)
      wg.Done()
			return
		}
//    fmt.Println("data sent")
	}
}
