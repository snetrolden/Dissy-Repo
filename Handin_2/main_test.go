package main

import (
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"testing"
	"time"
)

//We use blackbox testing throughout this hence some cases require good chunks of code for setup and proper testing
//Note we are too lazy to create helper methods right now.

func TestStartPeers(t *testing.T) {
	testPeer1 := newPeer(9999)
	testPeer2 := newPeer(8888)

	go testPeer1.Start()
	go testPeer2.Start()

	time.Sleep(1 * time.Second)

	conn1, err1 := net.Dial("tcp", "127.0.0.1:9999")
	if err1 != nil {
		t.Error("Could not conenct", err1)
	} else {
		conn1.Close()
	}

	conn2, err2 := net.Dial("tcp", "127.0.0.1:8888")
	if err2 != nil {
		t.Error("Could not connect", err2)
	} else {
		conn2.Close()
	}

}
func TestConnect(t *testing.T) {
	testPeer1 := newPeer(9991)
	testPeer2 := newPeer(8881)

	go testPeer1.Start()
	go testPeer2.Start()

	time.Sleep(1 * time.Second)

	err := testPeer1.Connect("127.0.0.1", 8881)
	if err != nil {
		t.Error("Peer 1 failed to connect to Peer 2:", err)
	}
	err = testPeer2.Connect("127.0.0.1", 9991)
	if err != nil {
		t.Error("Peer 2 failed to connect to Peer 1:", err)
	}

}

func TestSendKMessages(t *testing.T) {

	senderPort := 7000
	sender := newPeer(senderPort)
	go sender.Start()

	time.Sleep(1 * time.Second)

	receiverPort := 7001
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(receiverPort))
	if err != nil {
		t.Fatal(err) //fail that shi
	}
	defer listener.Close()

	err = sender.Connect("127.0.0.1", receiverPort)
	if err != nil {
		t.Fatal("Sender failed to connect:", err)
	}

	conn, err := listener.Accept()
	if err != nil {
		t.Fatal("Test listener failed to conenct:", err)
	}
	defer conn.Close()

	//handshake complete

	//loop 100 messages through the connection
	k := 100
	testMsg := &Message{
		Type:    "Test",
		MsgID:   "1",
		From:    "Sender",
		Payload: []byte("Hello"),
	}

	for i := 0; i < k; i++ {
		target := "127.0.0.1:" + strconv.Itoa(receiverPort)
		err := sender.Send(target, testMsg)
		if err != nil {
			t.Errorf("Failed to send message: %d %v", i, err)
		}
	}

	//we manually do what handleConnection do here for blackbox reasons and whatevs
	for i := 0; i < k; i++ {

		lenghtBuf := make([]byte, 4)
		_, err := io.ReadFull(conn, lenghtBuf)
		if err != nil {
			t.Fatalf("Failed to read length prefix for message %d: %v", i, err)
		}
		msgLenght := binary.BigEndian.Uint32(lenghtBuf)

		// Read Body
		msgBuf := make([]byte, msgLenght)
		_, err = io.ReadFull(conn, msgBuf)
		if err != nil {
			t.Fatalf("Failed to read message body for message %d: %v", i, err)
		}
	}

}
