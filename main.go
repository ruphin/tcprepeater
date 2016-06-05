package main

import "net"
import "fmt"
import "log"
import "io"
import DEATH "github.com/vrecan/death"
import SYS "syscall"

var death *DEATH.Death
var connections []io.Closer

func startTestListener() {
	test, err := net.Listen("tcp", "localhost:4001")
	connections = append(connections, test)
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
		connections = append(connections, conn)
		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		go func(c net.Conn, num int) {
			fmt.Println(num, " Starting")
			buf := make([]byte, 0, 4096)
			n := 0
			for {
				n += 1
				_, err := c.Read(buf)
				if err != nil && err != io.EOF {
					fmt.Println(num, " Something bad ", err)
					break
				}
				fmt.Println(num, " Read:", buf)
				if n == 10 {
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

func main() {
	death = DEATH.NewDeath(SYS.SIGINT, SYS.SIGTERM)
	connections = make([]io.Closer, 0)

	go startTestListener()
	var err error
	a, err = net.Dial("tcp", "localhost:4001")
	fmt.Println("Test1")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	b, err = net.Dial("tcp", "localhost:4001")
	fmt.Println("Test2")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	c, err = net.Dial("tcp", "localhost:4001")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	connections = append(connections, a)
	connections = append(connections, b)
	connections = append(connections, c)

	source, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		panic(err)
	}
	connections = append(connections, source)

	go func() {
		death.WaitForDeath(connections...)
	}()

	for {
		conn, err := source.Accept()
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		connections = append(connections, conn)
		go func(c net.Conn) {
			buf := make([]byte, 0, 4096)
			n := 0
			for {
				_, err := c.Read(buf)
				n += 1
				if err != nil && err != io.EOF {
					fmt.Println("Something bad:", err)
					break
				}
				send(buf)
				if n == 10 {
					break
				}
			}
			c.Close()
		}(conn)
	}

}
