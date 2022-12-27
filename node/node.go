package node

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/VANADAIN/drifter/dcrypto"
	"github.com/VANADAIN/drifter/settings"
	"github.com/VANADAIN/drifter/types"
)

type Node struct {
	id          *dcrypto.PublicKey
	ListenPort  string
	Lsn         net.Listener
	ConnLimiter int
	ActiveConns map[string]net.Conn // net.Addr -> connection
	ConnList    map[string]string   // strings of pubk -> net.Addr.String()
	Aliases     map[string]string   // name (alias) -> pubk string
	Msgch       chan types.Message
	Quitch      chan struct{}
}

func NewNode(port string) *Node {
	priv := dcrypto.GeneratePrivateKey()
	pub := priv.Public()

	settings := settings.Read()
	keys_path := settings.Keys_path

	checkKeysSaved(keys_path)
	fullpath := keys_path + "/key.private"

	saveKeyToFile(priv, fullpath)
	log.Printf("Private key of node saved to key.private")

	return &Node{
		id:         pub,
		ListenPort: port,
		Quitch:     make(chan struct{}),
		Msgch:      make(chan types.Message, 10),
	}
}

// load node with existing keys and settings
func LoadNode() *Node {
	priv := readKeyFromFile("./keys_saved/key.private")
	pub := priv.Public()

	return &Node{
		id:         pub,
		ListenPort: ":3000",
		Quitch:     make(chan struct{}),
		Msgch:      make(chan types.Message, 10),
	}
}

func (s *Node) Start() error {
	ln, err := net.Listen("tcp", s.ListenPort)
	if err != nil {
		return err
	}

	defer ln.Close()

	s.Lsn = ln
	go s.acceptLoop()

	<-s.Quitch
	close(s.Msgch)

	return nil
}

func (s *Node) acceptLoop() {
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
			s.addActiveConnection(conn)
			s.addConnectionToList(conn)

			go s.readLoop(conn)

		} else {
			conn.Write([]byte("All connection slots are busy"))
			conn.Close()
			log.Println("Max number of connections reached")
		}

	}
}

func (s *Node) readLoop(conn net.Conn) {
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

		msg := types.Message{
			From:    conn.RemoteAddr().String(),
			Payload: buf[:n],
		}

		// fmt.Println(msg)
		s.Msgch <- msg

		conn.Write([]byte("msg received"))
	}
}

func (s *Node) addActiveConnection(conn net.Conn) {
	address := conn.RemoteAddr().String()

	s.ActiveConns[address] = conn
}

func (s *Node) addConnectionToList(conn net.Conn) {
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
