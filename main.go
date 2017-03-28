package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"reflect"
	"syscall"
)

func bindToIf(conn net.Conn, interfaceName string) {
	ptrVal := reflect.ValueOf(conn)
	val := reflect.Indirect(ptrVal)
	//next line will get you the net.netFD
	fdmember := val.FieldByName("fd")
	val1 := reflect.Indirect(fdmember)
	netFdPtr := val1.FieldByName("sysfd")
	fd := int(netFdPtr.Int())
	//fd now has the actual fd for the socket
	err := syscall.SetsockoptString(fd, syscall.SOL_SOCKET,
		syscall.SO_BINDTODEVICE, interfaceName)
	if err != nil {
		log.Fatal(err)
	}
}

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
		local, err := net.ResolveUDPAddr("udp", ":5060")
		if err != nil {
			panic(err)
		}

		remote, err := net.ResolveUDPAddr("udp", "localhost:8888")
		if err != nil {
			panic(err)
		}

		conn, err := net.DialUDP("udp4", local, remote)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		bindToIf(conn, "bond0")

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
