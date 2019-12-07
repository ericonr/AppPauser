package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Requires argument:")
		fmt.Println("  pause")
		fmt.Println("  resume")
		fmt.Println("  kill")
		os.Exit(1)
	}

	conn, err := net.Dial("unix", "/tmp/app-pauser.socket")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Socket: ", conn.LocalAddr)

	arg := os.Args[1]
	conn.Write([]byte(arg))
}
