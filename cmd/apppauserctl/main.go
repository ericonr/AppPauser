package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ericonr/AppPauser/internal/apppauser"
)

var commands = apppauser.AvailableCommands

func print_help() {
	fmt.Println(os.Args[0], " requires argument:")
	for _, cmd := range commands {
		fmt.Println("  ", cmd)
	}
	os.Exit(1)
}

func main() {
	if len(os.Args) == 1 {
		print_help()
	}

	command := os.Args[1]

	valid_command := false
	for _, cmd := range commands {
		if command == cmd {
			valid_command = true
			break
		}
	}

	if !valid_command {
		print_help()
	}

	conn, err := net.Dial("unix", apppauser.SocketPath())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Socket: ", conn.LocalAddr)

	conn.Write([]byte(command))
	conn.Close()
}
