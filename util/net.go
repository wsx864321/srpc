package util

import "net"

// ReadFixData 读取网络连接对端发送的固定长度的内容
func ReadFixData(conn net.Conn, buf []byte) error {
	var (
		pos       = 0
		totalSize = len(buf)
	)
	for {
		c, err := conn.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos = pos + c
		if pos == totalSize {
			break
		}
	}

	return nil
}

// Write 对网络连接对端发送内容
func Write(conn net.Conn, data []byte) error {
	totalLen := len(data)
	writeLen := 0
	for {
		len, err := conn.Write(data[writeLen:])
		if err != nil {
			return err
		}
		writeLen = writeLen + len
		if writeLen >= totalLen {
			break
		}
	}
	return nil
}
