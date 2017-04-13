package main


func Send(message TcpMessage) {
	SendSocket(message)
}

func Receive() {
	for {
		ReceiveSocket()
		//fmt.Println(recvmsg)
	}
}
