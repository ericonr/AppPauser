package apppauser

import (
	"os"
	"path"
)

var AvailableCommands = []string{"pause", "resume", "toggle", "kill"}

func SocketPath() string {
	socket := os.Getenv("APPPAUSER_SOCK")
	if len(socket) != 0 {
		return socket
	}

	runtime := os.Getenv("XDG_RUNTIME_DIR")
	if len(runtime) == 0 {
		runtime = "/tmp"
	}

	qualifier := os.Getenv("APPPAUSER")

	return path.Join(runtime, "apppauser-"+qualifier+".socket")
}
