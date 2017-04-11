package main

import (
    "net"
    "bufio"
    "os"
    "fmt"
)

var connections = make(map[string]Connection)
var queue = make(chan TcpMessage)

type Connection struct {
    Name    string
    Address string
}


func connect(name, address string) {
    go listen(address)

    connections[name] = Connection{name, address}
}

func listen(address string) error {
    conn, error := net.Dial("tcp", address)
    for {
        message, error := bufio.NewReader(conn).ReadString('\n')
        if error != nil {
            break
        }
        queue <- fromJsonString(message)
    }
    return error
}

func sample() {

    // connect to this socket
    conn, _ := net.Dial("tcp", "127.0.0.1:8081")
    for {
        // read in input from stdin
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Text to send: ")
        text, _ := reader.ReadString('\n')
        // send to socket
        fmt.Fprintf(conn, text + "\n")
        // listen for reply
        message, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Print("UIMessage from server: " + message)
    }
}
