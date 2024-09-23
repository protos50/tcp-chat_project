package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var clients = make(map[net.Conn]string)
var mutex = &sync.Mutex{} // Para sincronización al modificar la lista de clientes

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Server error:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Server started at localhost:8080...")

	// Goroutine para permitir que el servidor envíe mensajes
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Message from the server: ")
			if scanner.Scan() {
				message := scanner.Text()
				broadcastMessage(nil, "[Server]: "+message+"\n")
			}
		}
	}()

	// Manejar conexiones entrantes
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Manejar cada conexión en una goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// add new client to list
	mutex.Lock()
	clients[conn] = conn.RemoteAddr().String()
	mutex.Unlock()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr().String())
	// Enviar la IP y puerto del cliente
	conn.Write([]byte(fmt.Sprintf("Conectado como: %s\n", conn.RemoteAddr().String())))

	// Read and broadcast messages
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())

			// Remove the client from the list
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			return
		}

		// /listar command
		if strings.TrimSpace(message) == "/listar" {
			fmt.Println("Comando /listar recibido")
			listUsers(conn)
		} else {
			// Broadcast message to all clients
			fmt.Printf("Message from %s: %s", conn.RemoteAddr().String(), message)
			broadcastMessage(conn, message)
		}
	}
}

func listUsers(conn net.Conn) {
	mutex.Lock()

	var userList []string
	for _, user := range clients {
		userList = append(userList, user)
	}

	mutex.Unlock()

	var message string
	if len(userList) == 0 {
		message = "No users connected.\n"
	} else {
		message = "Users connected:\n" + strings.Join(userList, "\n") + "\n"
	}

	// Enviar el mensaje al cliente solicitante
	conn.Write([]byte(message))

	// Imprimir en el servidor para depuración (opcional)
	fmt.Println("Clients list:")
	for _, user := range clients {
		fmt.Printf("Cliente: %s\n", user)
	}
}

func broadcastMessage(sender net.Conn, message string) {
	// Limpiar el mensaje
	cleanMessage := strings.TrimSpace(message)

	// Formato del mensaje para mostrar quién lo envió
	fullMessage := fmt.Sprintf("[%s]: %s\n", clients[sender], cleanMessage)

	// Enviar a todos los clientes, excepto al que envió el mensaje
	mutex.Lock()
	for client := range clients {
		if client != sender {
			client.Write([]byte(fullMessage))
		}
	}
	mutex.Unlock()

	fmt.Printf("Broadcast message: %s", fullMessage)
}
