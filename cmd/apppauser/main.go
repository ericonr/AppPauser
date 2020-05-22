package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/ericonr/AppPauser/internal/apppauser"
)

func CreateSocket() net.Listener {
	path := apppauser.SocketPath()
	conn, err := net.Listen("unix", path)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func CloseSocket(conn net.Listener) {
	conn.Close()
	log.Println("Closed socket: ", conn.Addr())
}

type ProcessStatus int

// Enums for the status of the running process
const (
	Paused  ProcessStatus = iota
	Running ProcessStatus = iota
)

func Pause(process *os.Process) ProcessStatus {
	log.Println("Pausing: ", process.Pid)
	err := process.Signal(syscall.SIGSTOP)
	if err != nil {
		log.Println("Unable to pause application.")
		return Running
	} else {
		return Paused
	}
}

func Resume(process *os.Process) ProcessStatus {
	log.Println("Resuming: ", process.Pid)
	err := process.Signal(syscall.SIGCONT)
	if err != nil {
		log.Println("Unable to resume application.")
		return Paused
	} else {
		return Running
	}
}

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
				status = Pause(process)
			} else {
				log.Println("Application is already paused.")
			}
		} else if clean_read == "resume" {
			if status == Running {
				log.Println("Application is already running.")
			} else {
				status = Resume(process)
			}
		} else if clean_read == "toggle" {
			if status == Running {
				status = Pause(process)
			} else if status == Paused {
				status = Resume(process)
			}
		} else if clean_read == "kill" {
			log.Println("Killing: ", process.Pid)
			if status == Paused {
				status = Resume(process)
			}
			process.Signal(syscall.SIGTERM)
			break
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
