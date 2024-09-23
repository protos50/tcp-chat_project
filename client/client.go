package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Server connection error:", err)
		return
	}
	defer conn.Close()

	// Lanzar goroutine para recibir mensajes del servidor
	go receiveMessages(conn)

	// Leer mensajes del usuario y enviarlos al servidor
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Start messaging: \n")
	for {
		if scanner.Scan() {
			message := scanner.Text()

			// Salir si el usuario escribe "/quitar"
			if message == "/quitar" {
				fmt.Println("Disconnected from server...")
				return
			}

			// Enviar el comando /listar al servidor
			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error sending message:", err)
				return
			}
		} else {
			fmt.Println("Invalid input.")
		}
	}
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		// Si el mensaje contiene la dirección del cliente, no imprimirlo
		if !strings.Contains(message, conn.RemoteAddr().String()) {
			fmt.Print("Server message: " + message)
		}
	}
}
