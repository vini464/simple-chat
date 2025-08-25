package main

import (
	"fmt"
	"net"
	"sync"
	"vini464/simple-chat/utils"
)

var CLIENTS = make(map[string]chan utils.Message)
var QUEUE   = make([]string, 0)
var ROOMS   = make(map[string]string)

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
  var username string
	var wg_clients sync.WaitGroup
	receive_channel := make(chan []byte)
	data_to_send := make(chan []byte)
  self_channel := make(chan utils.Message)

	defer conn.Close()
	defer wg_clients.Wait()
  defer delete(CLIENTS, username)
  

	fmt.Println("[log] Novo cliente:", conn.RemoteAddr().String())

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
			_, ok := ROOMS[username]
      if ok {
        msg := utils.Message{Cmd: "message", Data: "["+username+"]:"+ message.Data}
        CLIENTS[ROOMS[username]] <- msg
      } 
		case "set_user":
			var response utils.Message
			_, ok := CLIENTS[message.Data]
			if ok {
				response = utils.Message{Cmd: "error", Data: "invalid user"}
        sendResponse(response, data_to_send)
			} else {
				response = utils.Message{Cmd: "ok", Data: "user saved"}
				CLIENTS[message.Data] = self_channel
        username = message.Data
        sendResponse(response, data_to_send)
        if (len(QUEUE) == 0) {
          QUEUE = utils.Enqueue(QUEUE, username)
        } else {
          var other string
          other, QUEUE = utils.Dequeue(QUEUE)
          ROOMS[username] = other
          ROOMS[other] = username
          response := utils.Message{Cmd: "allocated", Data: username}
          CLIENTS[other] <-response
          response = utils.Message{Cmd: "allocated", Data: other}
          sendResponse(response, data_to_send)
        }
			}
		default:
			fmt.Println("[debug]: unknow command")
		}
	}
}

func sendResponse(msg utils.Message, send_channel chan []byte) {
			serialized, err := utils.SerializeJson(msg)
			for err != nil {
				fmt.Println("[error] - error while serializing:", err)
				serialized, err = utils.SerializeJson(msg)
			}
			send_channel <- serialized
}

