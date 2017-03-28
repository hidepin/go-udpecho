package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("Server is running at localhost:8888")
		conn, err := net.ListenPacket("udp", "localhost:8888")
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		buffer := make([]byte, 65535)
		for {
			length, remoteAddress, err := conn.ReadFrom(buffer)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Received from %v: %v\n",
				remoteAddress, string(buffer[:length]))
			_, err = conn.WriteTo([]byte("Hello from Server"), remoteAddress)
			if err != nil {
				panic(err)
			}
		}
	} else {
		conn, err := net.Dial("udp4", "localhost:8888")
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		fmt.Println("Sending to server")
		_, err = conn.Write([]byte("Hello from Client"))
		if err != nil {
			panic(err)
		}
		fmt.Println("Receiving from server")
		buffer := make([]byte, 65535)
		length, err := conn.Read(buffer)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Received: %s\n", string(buffer[:length]))
	}
}
