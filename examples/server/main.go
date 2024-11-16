package main

import (
	"github.com/nonya123456/cobble"
)

func main() {
	s := cobble.Server{Addr: ":25565"}
	s.Run()
}
