package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:4040")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ch := make(chan string)
	go read(conn, ch)
	go write(conn, ch)
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, "\n")
		var s string
		switch text {
		case "0":
			m := map[string]string{"cmd": "upload_score", "score": "A"}
			b, err := json.Marshal(m)
			if err != nil {
				fmt.Println("got error:", err.Error())
			}
			s = string(b)
		case "1":
			s = "[UR]100000000111111"
		case "2":
			m := map[string]string{"cmd": "confirm_btn"}
			b, _ := json.Marshal(m)
			s = string(b)
		case "3":
			m := map[string]string{"cmd": "confirm_init_score"}
			b, _ := json.Marshal(m)
			s = string(b)
		}
		ch <- s
	}
}

func read(conn net.Conn, ch chan string) {
	r := bufio.NewReader(conn)
	for {
		b, err := r.ReadByte()
		if err != nil {
			fmt.Println("err:", err.Error())
			os.Exit(1)
		}
		if b != 60 {
			fmt.Println("hardware message must start with <")
			os.Exit(1)
		}
		msg := make([]byte, 0)
		for {
			b, err := r.ReadByte()
			if err != nil {
				fmt.Println("err:", err.Error())
				os.Exit(1)
			}
			if b == 62 {
				break
			}
			msg = append(msg, b)
		}
		if len(msg) == 0 {
			fmt.Println("got empty hardware message")
			os.Exit(1)
		}
		msgStr := string(msg)
		fmt.Println("got message: ", msgStr)
	}
}

func write(conn net.Conn, ch chan string) {
	for {
		s := <-ch
		fmt.Println("will send:", s)
		fmt.Fprintf(conn, "<"+s+">")
	}
}
