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
	receive_channel := make(chan []byte)
	data_to_send := make(chan []byte)
	wg_clients.Add(1)
	go utils.ReceiveHandler(conn, receive_channel, &wg_clients)
	wg_clients.Add(1)
	go utils.SendHandler(conn, data_to_send, &wg_clients)
	for {
		income := <-receive_channel
		var message utils.Message
		utils.DeserializeToJson(income, &message)
		switch message.Cmd {
		case "message":
			fmt.Println("[debug]: received message:", message.Data)
			response := utils.Message{Cmd: "message", Data: "[server]: i received " + message.Data}
			serialized, err := utils.SerializeJson(response)
			if err != nil {
				fmt.Println("[error] - error while serializing:", err)
			} else {
				data_to_send <- serialized
			}

		default:
			fmt.Println("[debug]: unknow command")
		}
	}
	wg_clients.Wait()
}
