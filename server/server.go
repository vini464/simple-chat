package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "7070"
	SERVER_TYPE = "tcp"
	SERVER_PATH = SERVER_HOST + ":" + SERVER_PORT
)

func main() {

	fmt.Println("[log] Iniciando servidor...")
	listener, err := net.Listen(SERVER_TYPE, SERVER_PATH)
	if err != nil {
		fmt.Println("[error] Ocorreu um erro:", err.Error())
	}
	fmt.Println("[log] Escutando na porta:", SERVER_PORT)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[error] Não foi possível aceitar a conexão:", err.Error())
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("[log] Novo cliente:", conn.RemoteAddr().String())
	//	send_channel := make(chan string)
	receive_channel := make(chan string)
	data_to_send := make(chan string)

	//	go sender(conn, send_channel)
	go receiver(conn, receive_channel)
	go sender(conn, data_to_send)

	<-receive_channel

}

func sender(conn net.Conn, msg chan string) {
	for {
		data := <-msg
		size := uint32(len(data))
		header := make([]byte, 4)

		binary.BigEndian.PutUint32(header, size)
		_, err := conn.Write(header)
		if err != nil {
			fmt.Println("[error] - connection lost")
			return
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("[error] - connection lost")
			return
		}
	}
}

func receiver(conn net.Conn, receive_channel chan string) {
	header := make([]byte, 4)
	msg := ""
	// Loop de espera pela mensagem
	for {
		_, err := conn.Read(header)
		if err != nil {
			if err == io.EOF {
				fmt.Println("[log] O cliente encerrou a conexão")
			} else {
				fmt.Println("[error] O client foi desconectado:", err.Error())
			}
			receive_channel <- "finished"
			return // encerra a função caso algum erro tenha ocorrido
		}
		// loop de leitura da mensagem
		size := binary.BigEndian.Uint32(header)
		data := make([]byte, int(size))
		readed := 0
		fmt.Println("[debug] - expecting", int(size), "bytes")
		for strLen, err := conn.Read(data); readed < int(size) || err != nil; {
			fmt.Println("strLen: ", strLen)
			if err != nil {
				if err == io.EOF {
					fmt.Println("[log] O cliente encerrou a conexão")
				} else {
					fmt.Println("[error] O client foi desconectado:", err.Error())
				}
				receive_channel <- "finished"
				return // encerra a função caso algum erro tenha ocorrido
			}

			msg += string(data[:strLen])
			readed += strLen
			strLen = 0
		}
		if readed == int(size) && msg != "" {
			fmt.Println("[log] received message from client", conn.RemoteAddr().String(), ":", msg)
			msg = ""
		}
	}
	// TODO: create a simple protocol like: {method: "name"; data: "info"}
	// TODO: create a handler for each message
}
