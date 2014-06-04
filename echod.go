package main

import (
	"fmt"
	"io"
	"net"
)

func connection_handler(conn net.Conn) {
	fmt.Printf("Connection handler started for %s\n", conn.RemoteAddr())
	io.Copy(conn, conn)
	conn.Close()
	fmt.Printf("Connection closed by %s\n", conn.RemoteAddr())
}

func main() {

	srv, err := net.Listen("tcp", ":2323")
	if err != nil {
		fmt.Println("Error starting service")
		fmt.Printf("Error: %s", err)
	}
	defer srv.Close()

	for {

		conn, err := srv.Accept()
		if err != nil {
			fmt.Println("Error accepting incoming connection")
			fmt.Printf("Error: %s", err)
		}
		fmt.Printf("Incoming connection from %s on %s\n", conn.RemoteAddr(), conn.LocalAddr())
		go connection_handler(conn)
	}

}