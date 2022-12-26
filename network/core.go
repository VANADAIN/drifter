package network

import (
	"fmt"
	"io"
	"net"

	"github.com/VANADAIN/drifter/node"
)

func Start(s *node.Node) error {
	ln, err := net.Listen("tcp", s.ListenPort)
	if err != nil {
		return err
	}

	defer ln.Close()

	s.Lsn = ln
	go acceptLoop(s)

	<-s.Quitch
	close(s.Msgch)

	return nil

}

func acceptLoop(s *node.Node) {
	for {
		conn, err := s.Lsn.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("new conn: ", conn.RemoteAddr())

		// todo: add new connection, add active conn

		go readLoop(s, conn)
	}
}

func readLoop(s *node.Node, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("Connection closed")

				// todo: delete active conn

				break
			} else {
				fmt.Println("read error: ", err)
			}
		}

		msg := node.Message{
			From:    conn.RemoteAddr().String(),
			Payload: buf[:n],
		}

		// fmt.Println(msg)
		s.Msgch <- msg

		conn.Write([]byte("msg received"))
	}
}

func addActiveConnection(s *node.Node, conn net.Conn) {
	s.ActiveConns = append(s.ActiveConns, conn)
}
