package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
)

type Client struct {
	conn     net.Conn
	username string
	writer   *bufio.Writer
}

var (
	clients = make(map[net.Conn]*Client)
	mu      sync.Mutex
)

<<<<<<< HEAD
var NGROK_TOKEN = "YOUR TOKEN HERE"

=======
var YOUR_TOKEN = "YOUR_TOKEN_HERE"
>>>>>>> 825bfed (Fixed some bugs :D)

func main() {
	// Start ngrok tunnel
	if err := startNgrok(); err != nil {
		fmt.Println("Error starting ngrok:", err)
		return
	}

	// Start TCP server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Local server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func startNgrok() error {
	ctx := context.Background()
	listener, err := ngrok.Listen(ctx,
		config.TCPEndpoint(),
<<<<<<< HEAD
		ngrok.WithAuthtoken(NGROK_TOKEN),
=======
		ngrok.WithAuthtoken(YOUR_TOKEN),
>>>>>>> 825bfed (Fixed some bugs :D)
	)
	if err != nil {
		return err
	}

	fmt.Println("ngrok tunnel created:", listener.URL())

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting ngrok connection:", err)
				continue
			}
			localConn, err := net.Dial("tcp", "localhost:8080")
			if err != nil {
				fmt.Println("Error connecting to local server:", err)
				conn.Close()
				continue
			}
			go forward(localConn, conn)
			go forward(conn, localConn)
		}
	}()

	return nil
}

func forward(dst net.Conn, src net.Conn) {
	defer dst.Close()
	defer src.Close()
	bufio.NewReader(src).WriteTo(dst)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	writer := bufio.NewWriter(conn)
	writer.WriteString("Enter your username: ")
	writer.Flush()

	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	username := strings.TrimSpace(scanner.Text())

	client := &Client{conn: conn, username: username, writer: writer}

	mu.Lock()
	clients[conn] = client
	mu.Unlock()

	broadcast(fmt.Sprintf("%s has joined the chat", username), client)

	for {
		// Display prompt to the user
		writer.WriteString(fmt.Sprintf("[%s] > ", username))
		writer.Flush()

		if !scanner.Scan() {
			break
		}

		message := strings.TrimSpace(scanner.Text())

		if message != "" {
			broadcast(message, client) // Broadcasting only the message
		}
	}

	mu.Lock()
	delete(clients, conn)
	mu.Unlock()

	broadcast(fmt.Sprintf("%s has left the chat", username), client)
}

func broadcast(message string, sender *Client) {
	formattedMsg := fmt.Sprintf("[%s] %s", sender.username, message)

	mu.Lock()
	defer mu.Unlock()

	// Broadcast to all clients except the sender
	for _, client := range clients {
		if client != sender {
			// Clear the current prompt line before displaying the message
			client.writer.WriteString("\r") // Carriage return to clear prompt
			client.writer.WriteString("\033[K") // ANSI escape code to clear the line
			client.writer.WriteString(formattedMsg + "\n") // Broadcast the message
			client.writer.Flush()

			// Re-display prompt for each client after displaying the message
			client.writer.WriteString(fmt.Sprintf("[%s] > ", client.username)) // Display the prompt again
			client.writer.Flush()
		}
	}
}

