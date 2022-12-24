package main

import "github.com/VANADAIN/drifter/node"

// func (s *Server) Start() error {
// 	ln, err := net.Listen("tcp", s.listenAddr)
// 	if err != nil {
// 		return err
// 	}

// 	defer ln.Close()
// 	s.ln = ln

// 	go s.acceptLoop()

// 	<-s.quitch
// 	close(s.msgch)

// 	return nil
// }

// func (s *Server) acceptLoop() {
// 	for {
// 		conn, err := s.ln.Accept()
// 		if err != nil {
// 			fmt.Println("accept error: ", err)
// 			continue
// 		}

// 		fmt.Println("new conn: ", conn.RemoteAddr())

// 		go s.readLoop(conn)
// 	}

// }

// func (s *Server) readLoop(conn net.Conn) {

// 	defer conn.Close()

// 	buf := make([]byte, 2048)
// 	for {
// 		n, err := conn.Read(buf)
// 		if err != nil {
// 			if err == io.EOF {
// 				fmt.Println("EOF")
// 				return
// 			} else {
// 				fmt.Println("read error: ", err)
// 			}
// 		}

// 		msg := Message{
// 			from:    conn.RemoteAddr().String(),
// 			payload: buf[:n],
// 		}

// 		// fmt.Println(msg)
// 		s.msgch <- msg

// 		conn.Write([]byte("msg received"))
// 	}
// }

// type RouterClient interface {
// 	AddPortMapping(
// 		NewRemoteHost string,
// 		NewExternalPort uint16,
// 		NewProtocol string,
// 		NewInternalPort uint16,
// 		NewInternalClient string,
// 		NewEnabled bool,
// 		NewPortMappingDescription string,
// 		NewLeaseDuration uint32,
// 	) (err error)

// 	GetExternalIPAddress() (
// 		NewExternalIPAddress string,
// 		err error,
// 	)
// }

// func PickRouterClient(ctx context.Context) (RouterClient, error) {
// 	tasks, _ := errgroup.WithContext(ctx)
// 	// Request each type of client in parallel, and return what is found.
// 	var ip1Clients []*internetgateway2.WANIPConnection1
// 	tasks.Go(func() error {
// 		var err error
// 		ip1Clients, _, err = internetgateway2.NewWANIPConnection1Clients()
// 		return err
// 	})
// 	var ip2Clients []*internetgateway2.WANIPConnection2
// 	tasks.Go(func() error {
// 		var err error
// 		ip2Clients, _, err = internetgateway2.NewWANIPConnection2Clients()
// 		return err
// 	})
// 	var ppp1Clients []*internetgateway2.WANPPPConnection1
// 	tasks.Go(func() error {
// 		var err error
// 		ppp1Clients, _, err = internetgateway2.NewWANPPPConnection1Clients()
// 		return err
// 	})

// 	if err := tasks.Wait(); err != nil {
// 		return nil, err
// 	}

// 	// Trivial handling for where we find exactly one device to talk to, you
// 	// might want to provide more flexible handling than this if multiple
// 	// devices are found.
// 	switch {
// 	case len(ip2Clients) == 1:
// 		return ip2Clients[0], nil
// 	case len(ip1Clients) == 1:
// 		return ip1Clients[0], nil
// 	case len(ppp1Clients) == 1:
// 		return ppp1Clients[0], nil
// 	default:
// 		return nil, errors.New("multiple or no services found")
// 	}
// }

// func GetIPAndForwardPort(ctx context.Context) error {
// 	client, err := PickRouterClient(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	externalIP, err := client.GetExternalIPAddress()
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("Our external IP address is: ", externalIP)

// 	return client.AddPortMapping(
// 		"",
// 		// External port number to expose to Internet:
// 		3000,
// 		// Forward TCP (this could be "UDP" if we wanted that instead).
// 		"TCP",
// 		// Internal port number on the LAN to forward to.
// 		// Some routers might not support this being different to the external
// 		// port number.
// 		3000,
// 		// Internal address on the LAN we want to forward to.
// 		"192.168.0.198",
// 		// Enabled:
// 		true,
// 		// Informational description for the client requesting the port forwarding.
// 		"Drifter",
// 		// How long should the port forward last for in seconds.
// 		// If you want to keep it open for longer and potentially across router
// 		// resets, you might want to periodically request before this elapses.
// 		3600,
// 	)
// }

func main() {
	node.NewNode(":3000")
	// ctx := context.Background()
	// err := GetIPAndForwardPort(ctx)
	// if err == nil {
	// 	fmt.Println("port forwarding created!")
	// }

	// server := NewServer(":3000")

	// go func() {
	// 	for msg := range server.msgch {
	// 		fmt.Printf("receiver msg from (%s): %s\n", msg.from, string(msg.payload))
	// 	}
	// }()

	// log.Fatal(server.Start())
}
