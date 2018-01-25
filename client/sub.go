package client

import (
	"bulletsbox/public"
)

// Sub is a subscriber
type Sub struct {
	Conn *Client
	List map[string]bool
}

// NewSub return a new Sub representing the given names.
func NewSub(c *Client, name ...string) *Sub {
	s := &Sub{c, make(map[string]bool)}
	for _, i := range name {
		s.List[i] = true
	}
	return s
}

// Receive func
func (s *Sub) Receive() (string, []byte, error) {
	if err := s.Conn.cmd(nil, s, public.CmdReserve, nil); err != nil {
		return "", nil, err
	}
	return s.Conn.readRes(true)
}
