package server

import (
	"encoding/binary"
	"strconv"
	"fmt"
	"net"
	"bulletsbox/public"
)

// Server struct
type Server struct {
	l     *net.TCPListener
	store map[string]*Queue
}

// NewServer func
func NewServer(network, addr string) (*Server, error) {
	s := new(Server)
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	tcpListener, err := net.ListenTCP(network, tcpAddr)
	if err != nil {
		return nil, err
	}
	s.l = tcpListener
	s.store = make(map[string]*Queue)
	return s, nil
}

// Run a server
func (s *Server) Run() {
	defer s.l.Close()
	for {
		conn, err := s.l.AcceptTCP()
		if err != nil {
			fmt.Print(err.Error())
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *net.TCPConn) {
	var queueName string
	var watchList = make(map[string]bool)
	fmt.Println("CONNECT " + conn.RemoteAddr().String())
	fmt.Println("GO FUNC START")
	for {
		cmdCode := make([]byte, 1)
		if _, err := conn.Read(cmdCode); err != nil {
			break
		}
		fmt.Print("CODE ")
		fmt.Println(cmdCode[0])
		switch cmdCode[0] {
		case public.CmdUse:
			name, err := public.ParseName(conn)
			if err != nil {
				break
			}
			queueName = name
			if _, exists := s.store[queueName]; !exists {
				s.store[queueName] = NewQueue()
			}
			fmt.Println("USE " + queueName)
			conn.Write([]byte{public.ResSuccess})
		case public.CmdSend:
			score, err := public.ParseUint32(conn)
			if err != nil {
				break
			}
			delay, err := public.ParseUint32(conn)
			if err != nil {
				break
			}
			body, err := public.ParseBody(conn)
			if err != nil {
				break
			}
			fmt.Println("SEND " + strconv.FormatUint(uint64(score), 10) + " " + strconv.FormatUint(uint64(delay), 10) + " " + string(body))
			s.store[queueName].enQueue(NewItem(score, delay, body))
			conn.Write([]byte{public.ResSuccess})
		case public.CmdWatch:
			name, err := public.ParseName(conn)
			if err != nil {
				break
			}
			fmt.Println("WATCH " + name)
			watchList[name] = true
			fmt.Println(watchList)
			if _, exists := s.store[name]; !exists {
				conn.Write([]byte{public.ResFail})
			}else{
				s.store[name].list[conn.RemoteAddr().String()] = true
				conn.Write([]byte{public.ResSuccess})
			}
		case public.CmdIgnore:
			name, err := public.ParseName(conn)
			if err != nil {
				break
			}
			fmt.Println("IGNORE " + name)
			if _, exists := s.store[name]; exists {
				delete(watchList, name)
				delete(s.store[name].list, conn.RemoteAddr().String())
			}
			conn.Write([]byte{public.ResSuccess})
		case public.CmdReserve:
			for w := range watchList {
				if s.store[w].reserve.lock{
					conn.Write([]byte{public.ResSuccess})
					conn.Write([]byte{uint8(len(w))})
					conn.Write([]byte(w))
					bodyLen := make([]byte, 4)
					binary.BigEndian.PutUint32(bodyLen, uint32(len(s.store[w].reserve.item.body)))
					fmt.Println("RECEIVE " + w + " " + string(s.store[w].reserve.item.body))
					conn.Write(bodyLen)
					conn.Write(s.store[w].reserve.item.body)
					s.store[w].list[conn.RemoteAddr().String()] = false
					unlock:=true
					for _, bool := range s.store[w].list{
						if bool {
							unlock=false
						}
					}
					if unlock {
						s.store[w].reserve.lock = false
					}
				} else {
					conn.Write([]byte{public.ResFail})
				}
			}
		default:
			conn.Write([]byte{public.ResUnknowCmd})
		}
	}
	fmt.Println("CLOSE " + conn.RemoteAddr().String())
	fmt.Println(s.store[queueName].list)
	if q, exists := s.store[queueName]; exists {
		if _, exists := q.list[conn.RemoteAddr().String()]; exists {
			delete(q.list, conn.RemoteAddr().String())
		}
	}
	fmt.Println(s.store[queueName].list)
	fmt.Println("GO FUNC STOP")
}
