package main

import "net"
import "fmt"
import "log"
import "io"
import "bufio"
import "time"

func startTestListener() {
	test, err := net.Listen("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	n := 0
	for {
		// Wait for a connection.
		conn, err := test.Accept()
		n += 1
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn, num int) {
			fmt.Println(num, "Starting")
			var status string
			var err error
			r := bufio.NewReader(c)
			for {
				status, err = r.ReadString('\n')
				fmt.Print(num, status)
				if err != nil && err != io.EOF {
					fmt.Println(num, "Something bad", err)
					break
				}
			}
			c.Close()
		}(conn, n)
	}
}

var a, b, c net.Conn

func send(buf []byte) {
	var err error
	_, err = a.Write(buf)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	_, err = b.Write(buf)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	_, err = c.Write(buf)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func relay(msg string) {
	var err error
	_, err = fmt.Fprintf(a, msg)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	_, err = fmt.Fprintf(b, msg)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	_, err = fmt.Fprintf(c, msg)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func runRelay() {
	var err error
	a, err = net.Dial("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	b, err = net.Dial("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	c, err = net.Dial("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	source, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := source.Accept()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		go func(c net.Conn) {
			var msg string
			var err error
			r := bufio.NewReader(c)
			for {
				msg, err = r.ReadString('\n')
				fmt.Print("Get:", msg)
				relay(msg)
				if err != nil {
					fmt.Println("Something bad", err)
					break
				}
			}
			c.Close()
		}(conn)
	}
}

func main() {

	go startTestListener()
	go runRelay()
	
	time.Sleep(time.Second * 1)
	test, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Fprintf(test, "THIS IS A TEST\n")
	time.Sleep(time.Second * 1)

}
