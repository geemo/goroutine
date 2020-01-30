package node

import (
	"bufio"
	"net"
)

func NetworkSource(addr string) <-chan int {
	out := make(chan int)
	go func() {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}
		reader := bufio.NewReader(conn)
		source := ReaderSource(reader, -1)
		for v := range source {
			out <- v
		}
		close(out)
	}()
	return out
}

func NetworkSink(addr string, in <-chan int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		writer := bufio.NewWriter(conn)
		defer writer.Flush()
		WriterSink(writer, in)
	}()
}
