package main

import (
	"encoding/binary" //getting numbers to binary form
	"encoding/json"   //marshal and unmarshaling
	"fmt"             //input/output API
	"io"              //handles reading of messages
	"maps"
	"net"     //API for handling network related queries
	"strconv" //Converting thingies to strings
	"sync"    //Mutex
)

// Our peer structure
type Peer struct {
	port       int                 // Its own port
	peers      map[string]net.Conn // list over connected peers
	knownPeers map[string]bool     // Address book for known peers
	seenMsgs   map[string]bool     // list of uniquely seen msgIDs to prevent broadcast storm
	lock       sync.Mutex          // thread safety
	ledger     *Ledger
}

type Message struct {
	Type    string // ID for what messages purpose is, used in 'onMessage'
	MsgID   string // identifier
	From    string // sender
	Payload []byte // Actual message
}

// Makes a new peer struct with associated port and return its pointer
func newPeer(listenPort int) *Peer {

	return &Peer{
		port:       listenPort,
		peers:      make(map[string]net.Conn),
		knownPeers: make(map[string]bool),
		seenMsgs:   make(map[string]bool),
		ledger:     MakeLedger(),
	}
}

// Handles listening and accepting
// enables the peer to act as a server
func (p *Peer) Start() error {

	//listens for a TCP connection. Converted listenport (p.port) to a string
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(p.port))

	//Error check
	if err != nil {
		return err
	}

	fmt.Println("Listening for connections on port : ", p.port)

	//contineous listening
	for {

		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		//lets us know who the conneciton has been made with
		//fmt.Println(p.port, "Got a connection from", conn.RemoteAddr().String()) //fake ass port
		// using goroutines each connection is handled by its own process
		go p.handleConnection(conn)

	}
}

// Creates an outgoing TCP connection
// Enables to peer to conenct to another peer that is listening
func (p *Peer) Connect(addr string, port int) error {
	fullAddr := addr + ":" + strconv.Itoa(port) //combine address and port
	conn, err := net.Dial("tcp", fullAddr)      // 'Dial' an address/port that is listening
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Lock when modifying data in Peer to prevent crash
	p.lock.Lock()
	p.peers[fullAddr] = conn
	p.knownPeers[fullAddr] = true
	p.lock.Unlock()

	// fmt.Println(p.port, "Connected to", fullAddr)

	//put the connection on a goroutine
	go p.handleConnection(conn)
	p.getPeerList(fullAddr)

	//satisfy annoying return
	return nil
}

// Called when a message is received
func (p *Peer) OnMessage(msg *Message) {
	//use switch for the type of message recieved 'Ping' and 'Pong'
	p.lock.Lock()
	if p.seenMsgs[msg.MsgID] {
		p.lock.Unlock()
		return

	}
	p.seenMsgs[msg.MsgID] = true
	p.lock.Unlock()

	// fmt.Println(p.port, "Recieve", msg.Type, "from", msg.From, "ID :", msg.MsgID) //Print after deduplicating

	switch msg.Type {

	case "Ping":
		reply := &Message{
			Type:  "Pong",
			MsgID: msg.MsgID,
			From:  "127.0.0.1:" + strconv.Itoa(p.port),
		}

		p.Send(msg.From, reply)

		//Case for dynamically growing Peer network
	case "Join":

		p.lock.Lock()
		p.knownPeers[msg.From] = true
		p.lock.Unlock()

		p.FloodMessage(msg)

	case "Transaction":
		var tx Transaction
		json.Unmarshal(msg.Payload, &tx)

		//add check for already seen messages, avoid duplicate transactions?

		p.ledger.Transaction(&tx)
		// fmt.Println("Transaction complete", tx.From, "--->", tx.To, "in the amount:", tx.Amount)

		p.FloodMessage(msg)

	case "GetPeers":

		//check for reverse connection
		p.lock.Lock()
		_, connected := p.peers[msg.From]
		p.lock.Unlock()

		if !connected {
			host, portString, _ := net.SplitHostPort(msg.From)
			port, _ := strconv.Atoi(portString)
			p.Connect(host, port)
		}

		//safely copy map
		p.lock.Lock()
		known := make(map[string]bool)
		maps.Copy(known, p.knownPeers)
		p.lock.Unlock()

		payload, _ := json.Marshal(known)

		reply := &Message{
			Type:    "PeerList",
			MsgID:   "list-" + msg.From + "-" + strconv.Itoa(p.port),
			From:    "127.0.0.1:" + strconv.Itoa(p.port),
			Payload: payload,
		}
		p.Send(msg.From, reply) // Send it straight back!

	case "PeerList":
		var newPeers map[string]bool
		json.Unmarshal(msg.Payload, &newPeers)

		//Connect to all new peers (Consider code duplication in Case for "Join")
		for addr := range newPeers {
			p.lock.Lock()
			_, connected := p.peers[addr]
			p.lock.Unlock()

			isConnected := addr == "127.0.0.1:"+strconv.Itoa(p.port)

			if !connected && !isConnected {
				host, portString, _ := net.SplitHostPort(addr)
				port, _ := strconv.Atoi(portString)
				p.Connect(host, port)
			}
		}

		//Flood the list and let them join
		joinMsg := &Message{
			Type:  "Join",
			MsgID: "join-" + strconv.Itoa(p.port),
			From:  "127.0.0.1:" + strconv.Itoa(p.port),
		}
		p.FloodMessage(joinMsg)

	}

}

// The goroutine handling the connection between peers
func (p *Peer) handleConnection(conn net.Conn) {
	defer conn.Close()

	//loop
	for {
		lenghtBuffer := make([]byte, 4)
		_, err := io.ReadFull(conn, lenghtBuffer)
		if err != nil {
			return //close if error occure
		}

		msgLenght := binary.BigEndian.Uint32(lenghtBuffer) //convert back to data/binary
		msgBuffer := make([]byte, msgLenght)
		io.ReadFull(conn, msgBuffer)

		//umarshal the message with the lenght of the message
		var msg Message
		json.Unmarshal(msgBuffer, &msg)

		p.OnMessage(&msg)
	}

}

// Helper function for sending a message
func (p *Peer) Send(to string, msg *Message) error {

	//Find peer in maps
	p.lock.Lock()
	conn, ok := p.peers[to]
	p.lock.Unlock()
	//error handling
	if !ok {
		return fmt.Errorf("Peer not found : %v", to)
	}

	//marshal
	data, err := json.Marshal(msg)
	//error handling
	if err != nil {
		return err
	}
	//get lenght for prefixing data
	lenght := uint32(len(data)) //we use 32 for consistency, people in RC(RegneCentralen) recommended it

	p.lock.Lock()
	binary.Write(conn, binary.BigEndian, lenght) // I read that BigEndian is the default?
	conn.Write(data)                             //Consider catching errors and protect the write with locks? (ask TA)
	p.lock.Unlock()

	return nil

}

// Disseminates messages to all peers
func (p *Peer) FloodMessage(msg *Message) {

	// Send to neighbors
	p.lock.Lock()
	peers := make([]string, 0, len(p.peers))
	for addr := range p.peers {
		if addr != msg.From {
			peers = append(peers, addr)
		}
	}
	p.lock.Unlock()

	for _, addr := range peers {
		p.Send(addr, msg)
	}

}

// Used when making a transaction insteadof anually sending a "Transaction" message via send
func (p *Peer) FloodTransaction(tx *Transaction) {

	payload, err := json.Marshal(tx)
	if err != nil {
		fmt.Println("FloodTransaction Marshalling failed: ", err)
		return
	}

	msg := &Message{
		Type:    "Transaction",
		MsgID:   tx.ID,
		From:    "127.0.0.1:" + strconv.Itoa(p.port),
		Payload: payload,
	}

	//prevent looping back
	p.lock.Lock()
	p.seenMsgs[msg.MsgID] = true
	p.lock.Unlock()

	p.ledger.Transaction(tx)

	//Consider adding print statement?

	p.FloodMessage(msg)

}

// Helper function to get the list of known peers from specified peer
func (p *Peer) getPeerList(addr string) {

	reqMsg := &Message{
		Type:  "GetPeers",
		MsgID: "req-" + strconv.Itoa(p.port) + "-" + addr,
		From:  "127.0.0.1:" + strconv.Itoa(p.port),
	}
	p.Send(addr, reqMsg)

}
