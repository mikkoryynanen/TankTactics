package network

import "fmt"

type TCPServer struct {}

func (tcp *TCPServer) Connect() {
	fmt.Println("Making server with TCP")
}
