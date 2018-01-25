package client

import (
	"encoding/binary"
	"bulletsbox/public"
)

// Pub is Publisher
type Pub struct {
	Conn *Client
	Name string
}

// NewPub return a new Pub representing the given name.
func NewPub(c *Client, name string) *Pub {
	p := &Pub{c, name}
	return p
}

// Send func
func (p *Pub) Send(body []byte, score uint32, delay uint32) error {
	scoreBytes := make([]byte, 4)
	delayBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(scoreBytes, score)
	binary.BigEndian.PutUint32(delayBytes, delay)
	if err := p.Conn.cmd(p, nil, public.CmdSend, body, scoreBytes, delayBytes); err != nil {
		return err
	}
	_, _, err := p.Conn.readRes(false)
	if err != nil {
		return err
	}
	return nil
}
