package main

import (
	"./cache"
	"./httpserver"
	"./server"
)

func main() {
	ca := cache.New("inmemory")
	go server.New(ca).Listen()
	httpserver.New(ca).Listen()
}
