package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var serialPort *serial.Port

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	go http.ListenAndServe(":8080", nil)

	c := &serial.Config{
		Name: "COM5",
		Baud: 9600,
	}

	var err error

	serialPort, err = serial.OpenPort(c)
	if err != nil {
		log.Fatal("Erro ao abrir a porta serial:", err)
	}
	defer serialPort.Close()

	// Loop para ler dados do Arduino
	for {
		time.Sleep(1 * time.Second)
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade:", err)
		return
	}
	defer conn.Close()
	clients[conn] = true
	log.Println("Novo cliente conectado")

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Erro ao ler mensagem:", err)
			delete(clients, conn)
			break
		}

		log.Printf("Mensagem recebida: %s", msg)

		if string(msg) == "requestData" {
			readFromSerial(conn)
		}
	}
}

func readFromSerial(conn *websocket.Conn) {
	_, err := serialPort.Write([]byte("requestData\n"))
	if err != nil {
		log.Println("Erro ao enviar comando para o Arduino:", err)
		return
	}

	// time.Sleep(100 * time.Millisecond)

	buf := make([]byte, 64)
	n, err := serialPort.Read(buf)
	if err != nil {
		log.Println("Error ao ler da porta serial:", err)
		return
	}

	if n > 0 {
		data := string(buf[:n])
		fmt.Printf("Dados recebido: %s\n", data)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(data)); err != nil {
			log.Println(err)
			conn.Close()
		}
	}
}
