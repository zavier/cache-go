package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:12346")

	conn, _ := net.DialTCP("tcp", nil, addr)

	_, err := conn.Write([]byte("S7 9 testkeytestvalue"))
	if err != nil {
		log.Println(err)
		return
	}
	result, _ := ioutil.ReadAll(conn)
	fmt.Println(string(result))
}
