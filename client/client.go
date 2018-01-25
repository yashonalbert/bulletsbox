package client

import (
	"bytes"
	"encoding/binary"
	"net"
	"bulletsbox/public"
)

// Client struct
type Client struct {
	c     *net.TCPConn
	buf   *bytes.Buffer
	use   string
	watch map[string]bool
}

// NewClient func
func NewClient(network, addr string) (*Client, error) {
	c := new(Client)
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	c.c = tcpConn
	c.buf = new(bytes.Buffer)
	return c, nil
}

// Close Client
func (c *Client) Close() error {
	return c.c.Close()
}

func (c *Client) cmd(p *Pub, s *Sub, cmdCode uint8, body []byte, args ...[]byte) error {
	if err := c.assort(p, s); err != nil {
		return err
	}
	if body != nil {
		bodyLen := make([]byte, 4)
		binary.BigEndian.PutUint32(bodyLen, uint32(len(body)))
		args = append(args, bodyLen)
		args = append(args, body)
	}
	c.writeBytes(cmdCode, args...)
	c.c.Write(c.buf.Bytes())
	c.buf.Reset()
	return nil
}

func (c *Client) assort(p *Pub, s *Sub) error {
	if p != nil && c.use != p.Name {
		if err := public.CheckName(p.Name); err != nil {
			return err
		}
		c.writeBytes(public.CmdUse, []byte{uint8(len([]byte(p.Name)))}, []byte(p.Name))
		c.use = p.Name
		c.c.Write(c.buf.Bytes())
		c.buf.Reset()
		if _,_,err := c.readRes(false); err != nil {
			return err
		}
	}
	if s != nil {
		for l := range s.List {
			if !c.watch[l] {
				if err := public.CheckName(l); err != nil {
					return err
				}
				c.writeBytes(public.CmdWatch, []byte{uint8(len(l))}, []byte(l))
				c.c.Write(c.buf.Bytes())
				c.buf.Reset()
				if _,_,err := c.readRes(false); err != nil {
					return err
				}
			}
		}
		for w := range c.watch {
			if !s.List[w] {
				c.writeBytes(public.CmdIgnore, []byte{uint8(len(w))}, []byte(w))
				c.c.Write(c.buf.Bytes())
				c.buf.Reset()
				if _,_,err := c.readRes(false); err != nil {
					return err
				}
			}
		}
		c.watch = make(map[string]bool)
		for l := range s.List {
			c.watch[l] = true
		}
	}
	return nil

}

func (c *Client) writeBytes(cmdCode uint8, args ...[]byte) {
	c.buf.WriteByte(cmdCode)
	for _, i := range args {
		c.buf.Write(i)
	}
}

func (c *Client) readRes(readData bool) (string, []byte, error) {
	readBytes := make([]byte, 1)
	if _, err := c.c.Read(readBytes); err != nil {
		return "", nil, err
	}
	if readBytes[0] != public.ResSuccess {
		return "", nil, public.ParseResError(readBytes[0])
	}
	if readData {
		name, err := public.ParseName(c.c)
		if err != nil {
			return "", nil, err
		}
		body, err := public.ParseBody(c.c)
		if err != nil {
			return "", nil, err
		}
		return name, body, nil
	}
	return "", nil, nil
}
