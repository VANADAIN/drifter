package network

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/VANADAIN/drifter/types"
)

type NetworkManager struct {
	info        types.NodeInfo
	ConnLimiter int
	ListenPort  string
	Lsn         net.Listener
	ConnList    map[string]string    // strings of pubk -> net.Addr(ip)
	ActiveConns map[string]net.Conn  // net.Addr(ip) -> connection
	Mediators   map[string]*Mediator // net.Addr(ip) -> message manager
	Aliases     map[string]string    // name (alias) -> pubk string
	Msgch       chan types.Message
	Quitch      chan struct{}
}

func (s *NetworkManager) Start() error {
	ln, err := net.Listen("tcp", s.ListenPort)
	if err != nil {
		return err
	}

	defer ln.Close()

	go s.acceptLoop()

	// connect to other nodes from this node
	s.CreateConnections()

	s.Lsn = ln

	<-s.Quitch
	close(s.Msgch)

	return nil
}

// this loop is for "reserved" connections to empty slots
// num of reserved = connection_limiter - active connections from createconnections()
// also peers can connect after active peers disconnect
func (s *NetworkManager) acceptLoop() {
	for {
		conn, err := s.Lsn.Accept()
		if err != nil {
			fmt.Println("Accept connection error: ", err)
			continue
		}

		fmt.Println("New conn: ", conn.RemoteAddr())

		if len(s.ActiveConns) < s.ConnLimiter {
			conn.Write([]byte("This peer has slots, running mediator for you... Send node info or you will be disconnected."))
			s.RunMediator(conn)

			// remote node will send nodeinfo after 2 sec
			// and mediator will add it to list
			// waiting 5 sec
			time.Sleep(5 * time.Second)
			s.ValidateConnection(conn)
		} else {
			conn.Write([]byte("All connection slots of this peer are busy. Closing connection."))
			conn.Close()
			log.Println("Max number of connections reached")
		}

	}
}

func (s *NetworkManager) CreateConnections() {

	if len(s.ConnList) == 0 {
		panic("No addresses in the list! Connection list is empty.")
	}

	// TODO: run this in goroutines with channel (faster connection)

	for _, address := range s.ConnList {
		if len(s.ActiveConns) == s.ConnLimiter {
			break
		}

		// remote node will wait in it's accept loop for incoming connections
		conn, err := net.Dial("tcp", address)
		if err != nil {
			fmt.Errorf("Can't connect to: %s", address)
			continue
		}

		// there we need to send our node info
		// remote node is validating incoming connection
		new_mediator := &Mediator{
			nm: s,
			c:  conn,
		}
		s.Mediators[conn.RemoteAddr().String()] = new_mediator

		new_mediator.RunReadLoop()

		// sleep 2 seconds before sending info (remote node is preparing it's mediator)
		time.Sleep(2 * time.Second)
		new_mediator.SendNodeInfo()
	}
}

func (s *NetworkManager) initialConnection(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic("Remote peer from provided address if offline")
	}

	// after that connection remote peer will ask for node info
	// that info must be available after first init already

	// create mediator
	new_mediator := &Mediator{
		nm: s,
		c:  conn,
	}

	s.Mediators[conn.RemoteAddr().String()] = new_mediator

	time.Sleep(2 * time.Second)
	new_mediator.RunReadLoop()
	new_mediator.SendNodeInfo()

	return conn
}

func (s *NetworkManager) RunMediator(conn net.Conn) {
	new_mediator := &Mediator{
		nm: s,
		c:  conn,
	}

	s.Mediators[conn.RemoteAddr().String()] = new_mediator
	new_mediator.RunReadLoop()
}

func (s *NetworkManager) AddActiveConnection(conn net.Conn) {
	_, ok := s.ActiveConns[conn.RemoteAddr().String()]
	if ok {
		log.Print("Error adding active connection. Already active")
		return
	} else {
		address := conn.RemoteAddr().String()
		s.ActiveConns[address] = conn
	}
}

// calling this from inner mediator
func (s *NetworkManager) AddConnectionToList(conn net.Conn, pub *dcrypto.PublicKey) {
	for _, v := range s.ConnList {
		if v == conn.RemoteAddr().String() {
			log.Println("Error adding connection, connection already in list")
			return
		} else {
			s.ConnList[pub.String()] = conn.RemoteAddr().String()
		}
	}
}

func (s *NetworkManager) DeleteActiveConnection(conn net.Conn) {
	address := conn.RemoteAddr().String()
	delete(s.ActiveConns, address)
}

func (s *NetworkManager) ValidateConnection(conn net.Conn) {
	// check if connection was added to active connections
	_, ok := s.ActiveConns[conn.RemoteAddr().String()]
	// if this connection exists (mediator added it)
	if ok {
		return
	} else {
		conn.Write([]byte("Your peer did not send node info. Closing connection..."))
		// delete mediator
		delete(s.Mediators, conn.RemoteAddr().String())
		conn.Close()
	}
}
