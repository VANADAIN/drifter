package network

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/VANADAIN/drifter/types"
)

type NetworkManager struct {
	info        types.NodeInfo
	ConnLimiter int
	ListenPort  string
	Lsn         net.Listener
	ConnList    map[string]string   // strings of pubk -> net.Addr(ip)
	ActiveConns map[string]net.Conn // net.Addr(ip) -> connection
	Mediators   map[string]Mediator // net.Addr(ip) -> message manager
	Aliases     map[string]string   // name (alias) -> pubk string
	Msgch       chan types.Message
	Quitch      chan struct{}
}

func (s *NetworkManager) Start() error {
	ln, err := net.Listen("tcp", s.ListenPort)
	if err != nil {
		return err
	}

	defer ln.Close()

	// connect to other nodes from this node
	s.CreateConnections()

	s.Lsn = ln
	go s.acceptLoop()

	<-s.Quitch
	close(s.Msgch)

	return nil
}

// this loop is for "reserved" connections to empty slots
// num of reserved = connection_limiter - active connections from createconnections()
func (s *NetworkManager) acceptLoop() {
	for {
		conn, err := s.Lsn.Accept()
		if err != nil {
			fmt.Println("Accept connection error: ", err)
			continue
		}

		fmt.Println("New conn: ", conn.RemoteAddr())

		if len(s.ActiveConns) < s.ConnLimiter {

			// check if connection is from other node
			// validateConnection
			s.AddActiveConnection(conn)
			s.AddConnectionToList(conn)

			go s.runMediator(conn)

		} else {
			conn.Write([]byte("All connection slots are busy"))
			conn.Close()
			log.Println("Max number of connections reached")
		}

	}
}

func (s *NetworkManager) runMediator(conn net.Conn) {

	// TODO: change this to mediator instance and running read loop in mediator

	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println("Connection closed")
				s.DeleteActiveConnection(conn)
				break
			} else {
				fmt.Println("read error: ", err)
			}
		}

		msgb := types.MessageBody{
			From:    conn.RemoteAddr().String(),
			Payload: buf[:n],
		}

		msg := types.Message{
			Body: msgb,
		}
		msg.AddHash()

		// fmt.Println(msg)

		s.Msgch <- msg

		conn.Write([]byte("msg received"))
	}
}

// this method is for connecting to "public" (always active) node
func (s *NetworkManager) iniitialConenction(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic("Remote peer from provided address if offline")
	}

	s.AddActiveConnection(conn)
	// after that connection remote peer will ask for node info
	// that info must be available after first init already
	return conn
}

func (s *NetworkManager) CreateConnections() {
	// loop through connection list and try to connect to peers
}

func (s *NetworkManager) DeleteActiveConnection(conn net.Conn) {
	address := conn.RemoteAddr().String()

	delete(s.ActiveConns, address)
}

func (s *NetworkManager) AddActiveConnection(conn net.Conn) {
	address := conn.RemoteAddr().String()

	s.ActiveConns[address] = conn
}

func (s *NetworkManager) AddConnectionToList(conn net.Conn) {
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

func (s *NetworkManager) validateConnection(conn net.Conn) {
	// create mediator and send some info
}
