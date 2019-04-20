package pipeline

import (
	"bufio"
	"net"
)

func NetWorkSink(addr string, in <-chan int) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		defer listener.Close()

		//for {} 一次性的pipeline 因此不用for accept
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

func NetworkSource(addr string) <-chan int{
	out := make(chan int)
	go func() {
		conn, err:= net.Dial("tcp", addr)
		if err != nil {
			panic(err)
		}

		r := ReaderSource(bufio.NewReader(conn), -1)
		for v:= range r{
			out <- v
		}
		close(out)
	}()
	return out
}