package main

import (
	"encoding/binary" //getting numbers to binary form
	"encoding/json"   //marshal and unmarshaling
	"fmt"             //input/output API
	"io"              //handles reading of messages
	"net"             //API for handling network related queries
	"strconv"         //Converting thingies to strings
	"sync"            //Mutex
)

// Our peer structure
type Peer struct {
	port       int                 // Its own port
	peers      map[string]net.Conn // list over connected peers
	knownPeers map[string]bool     // Address book for known peers
	lock       sync.Mutex          // thread safety

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
		port:  listenPort,
		peers: make(map[string]net.Conn),
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
		fmt.Println(p.port, "Got a connection from", conn.RemoteAddr().String())
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

	fmt.Println(p.port, "Connected to", fullAddr)

	//put the connection on a goroutine
	go p.handleConnection(conn)

	//satisfy annoying return
	return nil
}

// Called when a message is received
func (p *Peer) OnMessage(from string, msg *Message) {
	//use switch for the type of message recieved 'Ping' and 'Pong'

	fmt.Println(p.port, "Recieve", msg.Type, "from", msg.From, "ID :", msg.MsgID)

	//Case: if ping then pong
	if msg.Type == "Ping" {
		reply := &Message{
			Type:  "Pong",
			MsgID: msg.MsgID,
			From:  "127.0.0.1:" + strconv.Itoa(p.port),
		}

		p.Send(msg.From, reply)
	}

	if msg.Type == "Addr Book" {
		var connectedPeers map[string]bool
		json.Unmarshal(msg.Payload, &connectedPeers) //Unmarshal the []byte into the map

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

		p.OnMessage(conn.RemoteAddr().String(), &msg)
	}

}

// Helper function for sending a message
func (p *Peer) Send(to string, msg *Message) error {
	//Find peer in maps
	//remember to lock when interacting with Peer
	p.lock.Lock()
	conn, ok := p.peers[to]
	p.lock.Unlock()

	//error handling
	if !ok {
		return fmt.Errorf("Peer not found : %v", to) // return needs type error so Println is kaput
	}

	//marshal
	data, err := json.Marshal(msg)
	//error handling
	if err != nil {
		return err
	}
	//get lenght for prefixing data
	lenght := uint32(len(data)) //we use 32 for consistency, people in RC(RegneCentralen) recommended it

	//Convert the prefix, and send it before sending the data.
	binary.Write(conn, binary.BigEndian, lenght) // I read that BigEndian is the default?
	conn.Write(data)                             //Consider catching errors and protect the write with locks? (ask TA)

	return nil

}
