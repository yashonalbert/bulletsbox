package public

import (
	"encoding/binary"
	"net"
)

// ParseName func
func ParseName(conn *net.TCPConn) (string, error) {
	readBytes := make([]byte, 1)
	if _, err := conn.Read(readBytes); err != nil {
		return "", err
	}
	readBytes = make([]byte, readBytes[0])
	if _, err := conn.Read(readBytes); err != nil {
		return "", err
	}
	return string(readBytes), nil
}

// ParseUint32 func
func ParseUint32(conn *net.TCPConn) (uint32, error) {
	readBytes := make([]byte, 4)
	if _, err := conn.Read(readBytes); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(readBytes), nil
}

// ParseBody func
func ParseBody(conn *net.TCPConn) ([]byte, error) {
	readBytes := make([]byte, 4)
	if _, err := conn.Read(readBytes); err != nil {
		return nil, err
	}
	readBytes = make([]byte, binary.BigEndian.Uint32(readBytes))
	if _, err := conn.Read(readBytes); err != nil {
		return nil, err
	}
	return readBytes, nil
}
