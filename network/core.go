package network

import (
	"fmt"
	"io"
	"log"
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
			fmt.Println("Accept connection error: ", err)
			continue
		}

		fmt.Println("New conn: ", conn.RemoteAddr())

		if len(s.ActiveConns) < s.ConnLimiter {
			addActiveConnection(s, conn)
		} else {
			conn.Write([]byte("All connection slots are busy"))
			conn.Close()
			log.Println("Max number of connections reached")
		}

		addConnectionToList(s, conn)

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
	address := conn.RemoteAddr().String()

	s.ActiveConns[address] = conn
}

func addConnectionToList(s *node.Node, conn net.Conn) {
	// search by address (value)
	for _, v := range s.ConnList {
		if v == conn.RemoteAddr().String() {
			return
		} else {
			// connection not in list ->
			// get public key of remote node
			// s.ConnList[] = conn.RemoteAddr().String()
		}
	}
}
