package communication

import (
	"fmt"
	"net"
  "encoding/binary"
)

func handleReceive(conn net.Conn, received_data chan string) {
	header := make([]byte, 4)
	msg := ""
	for {
		_, err := conn.Read(header)
		if err != nil {
			fmt.Println("[error] - algo deu errado")
			return
		}
		size := binary.BigEndian.Uint32(header)
		data := make([]byte, size)
		readed := 0
		for strLen, err := conn.Read(data); err != nil || readed < int(size); {
			if err != nil {
				fmt.Println("[error] - connection lost")
				return
			}
			msg += string(data)
			readed += strLen
			strLen = 0
		}
		received_data <- msg
	}
}

func handleSend(conn net.Conn, msg chan string) {
	for {
		data := <-msg
		size := uint32(len(data))
		header := make([]byte, 4)

		binary.BigEndian.PutUint32(header, size)
		_, err := conn.Write(header)
		if err != nil {
			fmt.Println("[error] - connection lost")
			panic(err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("[error] - connection lost")
			panic(err)
		}
	}
}
