package main

import (
	"./cache"
	"./server"
)

func main() {
	c := cache.New("inmemory")
	server.New(c).Listen()
}
