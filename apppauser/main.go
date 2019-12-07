package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func CreateSocket() net.Listener {
	path := "/tmp/app-pauser.socket"
	conn, err := net.Listen("unix", path)
	if err != nil {
		log.Println("Removed: ", path)
		os.Remove(path)
		conn = CreateSocket()
	}
	return conn
}

func CloseSocket(conn net.Listener) {
	conn.Close()
	log.Println("Closed socket: ", conn.Addr())
}

const (
	Paused  = iota
	Running = iota
)

func SocketMonitor(listen net.Listener, process *os.Process) {
	status := Running

	for true {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Socket connected!")

		read := make([]byte, 256)
		conn.Read(read)

		clean_read := strings.TrimRight(string(read), string(0))
		if clean_read == "pause" {
			if status == Running {
				log.Println("Pausing: ", process.Pid)
				process.Signal(syscall.SIGSTOP)
				status = Paused
			} else {
				log.Println("Application is already paused.")
			}
		} else if clean_read == "resume" {
			if status == Running {
				log.Println("Application is already running.")
			} else {
				log.Println("Resuming: ", process.Pid)
				process.Signal(syscall.SIGCONT)
				status = Running
			}
		} else if clean_read == "kill" {
			log.Println("Killing: ", process.Pid)
			if status == Paused {
				process.Signal(syscall.SIGCONT)
			}
			process.Signal(syscall.SIGTERM)
		}
	}
}

func StartProcess() (*exec.Cmd, string) {
	if len(os.Args) == 1 {
		fmt.Println("Remember to include the command to be run:")
		fmt.Println("  apppauser command [arguments]")
		os.Exit(1)
	}

	program := os.Args[1]
	args := os.Args[2:]

	cmd := exec.Command(program, args...)

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running: ", program)
	log.Println("Args: ", args)
	log.Println("PID: ", cmd.Process.Pid)

	return cmd, program
}

func main() {
	cmd, program := StartProcess()

	conn := CreateSocket()
	go SocketMonitor(conn, cmd.Process)

	cmd.Wait()
	log.Println("Finished execution: ", program)

	CloseSocket(conn)
}
