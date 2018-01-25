package main

import (
	"bulletsbox/server"
	"fmt"
)

func main() {
	s, err := server.NewServer("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Print(err.Error())
	}
	s.Run()
}
