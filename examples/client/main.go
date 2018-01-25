package main

import (
	"bulletsbox/client"
	"fmt"
)

func main() {
	c, err := client.NewClient("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Println(err.Error())
	}
	p := client.NewPub(c, "default")
	if err := p.Send([]byte("abc"), 1024, 0); err != nil {
		fmt.Println(err.Error())
	}
	if err := p.Send([]byte("abc"), 0, 0); err != nil {
		fmt.Println(err.Error())
	}
	if err := p.Send([]byte("abc"), 2048, 0); err != nil {
		fmt.Println(err.Error())
	}
	if err := p.Send([]byte("abc"), 1024, 0); err != nil {
		fmt.Println(err.Error())
	}

    s := client.NewSub(c, "default")
	for {
		name, body, err := s.Receive()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("RECEIVE " + name + " " + string(body))
	}
}
