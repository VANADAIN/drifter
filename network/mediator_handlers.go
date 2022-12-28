package network

func registerHandler() {
	// check len of active connections
	// if len is full -> some connection was activated before this
	// no need to close connection, cuz connection will be closed by validateConnection()

	// if ok -> call addtolist
}
