package main

import (
	"fmt"
	"net"
	"sync"
	"vini464/simple-chat/utils"
)

const (
	SERVER_HOST = "server"
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
  fmt.Println("[log] Escutando em:", listener.Addr())
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
	var wg_clients sync.WaitGroup
	fmt.Println("[log] Novo cliente:", conn.RemoteAddr().String())
	//	send_channel := make(chan string)
	receive_channel := make(chan []byte)
	data_to_send := make(chan []byte)
	wg_clients.Add(1)
	go utils.ReceiveHandler(conn, receive_channel, &wg_clients)
	wg_clients.Add(1)
	go utils.SendHandler(conn, data_to_send, &wg_clients)
	for {
		msg := <-receive_channel
		fmt.Println("[debug] - received:", string(msg))
    response := "[server]: i received: " + string(msg)
		data_to_send <- []byte(response)
	}
	wg_clients.Wait()
}

// TODO: create a simple protocol like: {method: "name"; data: "info"}
// TODO: create a handler for each message
