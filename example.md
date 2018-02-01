# Example

Server start:

```go
s, err := server.NewServer("tcp", "127.0.0.1:3000")
	if err != nil {
		fmt.Print(err.Error())
	}
	s.Run()
```

Publisher jobs:

```go
c, err := client.NewClient("tcp", "127.0.0.1:3000")
if err != nil {
    fmt.Println(err.Error())
}
p := client.NewPub(c, "default")
if err := p.Send([]byte("abc"), 1024, 0); err != nil {
    fmt.Println(err.Error())
}
```

Subscriber jobs:

```go
c, err := client.NewClient("tcp", "127.0.0.1:3000")
if err != nil {
    fmt.Println(err.Error())
}
s := client.NewSub(c, "default")
name, body, err := s.Receive()
if err != nil {
    fmt.Println(err.Error())
}
```
