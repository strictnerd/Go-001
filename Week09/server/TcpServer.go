package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	//c := make(chan net.Conn, 2)
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go tcpHandler(tcpConn)
	}
}

func tcpHandler(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println("Closed:", ipStr)
		conn.Close()
	}()

	readChan := make(chan string)
	writeChan := make(chan string)
	//防止通道关闭
	stopChan := make(chan bool)
	//读操作
	go readConn(conn, readChan, stopChan)
	//写操作
	go writeConn(conn, writeChan, stopChan)

	for {
		select {
		case readStr := <-readChan:
			upper := strings.ToUpper(readStr)
			writeChan <- upper
		case stop := <-stopChan:
			if stop {
				break
			}
		}
	}
}

func readConn(conn net.Conn, readChan chan<- string, stopChan chan<- bool) {
	for {
		data := make([]byte, 1024)
		_, err := conn.Read(data)
		if err != nil {
			fmt.Println(err)
			break
		}

		strData := string(data)
		fmt.Println("Received:", strData)

		readChan <- strData
	}
	stopChan <- true
}

func writeConn(conn net.Conn, writeChan <-chan string, stopChan chan<- bool) {
	for {
		strData := <-writeChan
		_, err := conn.Write([]byte(strData))
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Println("Send:", strData)
	}

	stopChan <- true
}

/*

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println("Closed:", ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	i := 0
	for {
		message, err := reader.ReadString('\n')
		if err!=nil || err == io.EOF {
			break
		}
		fmt.Println(string(message))
		time.Sleep(time.Second*3)

		msg := time.Now().String() + conn.RemoteAddr().String() + "Server Say hello! \n"

		b := []byte(msg)

		conn.Write(b)

		i++
		if i > 10 {
			break
		}
	}
}
*/
